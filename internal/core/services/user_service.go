package services

import (
	"context"
	"errors"
	"log"
	"time"

	"calespero/internal/core/domain"
	"calespero/internal/core/ports"
	"calespero/pkg/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo   ports.UserRepository
	jwt        *auth.JWTManager
	workerPool chan struct{}
	resultChan chan error
	userChan   chan *domain.User
	tokenChan  chan string
}

func NewUserService(userRepo ports.UserRepository, jwt *auth.JWTManager) ports.UserService {
	service := &userService{
		userRepo:   userRepo,
		jwt:        jwt,
		workerPool: make(chan struct{}, 10), // Limit to 10 concurrent operations
		resultChan: make(chan error),
		userChan:   make(chan *domain.User),
		tokenChan:  make(chan string),
	}

	return service
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	// Get a worker from the pool
	select {
	case s.workerPool <- struct{}{}: // Acquire worker
		defer func() { <-s.workerPool }() // Release worker
	case <-ctx.Done():
		return ctx.Err()
	}

	// Start goroutine for user creation
	go func() {
		// Check if user already exists
		existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
		if err != nil {
			s.resultChan <- err
			return
		}
		if existingUser != nil {
			s.resultChan <- errors.New("user already exists")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			s.resultChan <- err
			return
		}

		// Create new user
		user.ID = uuid.New().String()
		user.Password = string(hashedPassword)
		user.CreatedAt = time.Now()
		user.LastConnection = time.Now()

		s.resultChan <- s.userRepo.CreateUser(ctx, user)
	}()

	// Wait for result or context cancellation
	select {
	case err := <-s.resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	// Get a worker from the pool
	select {
	case s.workerPool <- struct{}{}: // Acquire worker
		defer func() { <-s.workerPool }() // Release worker
	case <-ctx.Done():
		return "", ctx.Err()
	}

	// Start authentication goroutine
	go func() {
		user, err := s.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			s.resultChan <- err
			s.tokenChan <- ""
			return
		}
		if user == nil {
			s.resultChan <- errors.New("user not found")
			s.tokenChan <- ""
			return
		}

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			s.resultChan <- errors.New("invalid password")
			s.tokenChan <- ""
			return
		}

		// Update last connection in a separate goroutine
		go func() {
			if updateErr := s.userRepo.UpdateLastConnection(ctx, user.ID); updateErr != nil {
				log.Printf("Error updating last connection: %v", updateErr)
			}
		}()

		// Generate token
		token, err := s.jwt.GenerateToken(user.ID)
		if err != nil {
			s.resultChan <- err
			s.tokenChan <- ""
			return
		}

		s.resultChan <- nil
		s.tokenChan <- token
	}()

	// Wait for result or context cancellation
	select {
	case err := <-s.resultChan:
		if err != nil {
			return "", err
		}
		return <-s.tokenChan, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *userService) ValidateToken(token string) (string, error) {
	return s.jwt.ValidateToken(token)
}
