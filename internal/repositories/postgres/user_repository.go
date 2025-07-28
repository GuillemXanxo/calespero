package postgres

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"calespero/internal/core/domain"
	"calespero/internal/core/ports"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db    *sql.DB
	mutex sync.RWMutex
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `
		INSERT INTO users (id, email, password, phone, created_at, last_connection)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Phone,
		time.Now(),
		time.Now(),
	)

	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query := `
		SELECT id, email, password, phone, created_at, last_connection
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
		&user.LastConnection,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateLastConnection(ctx context.Context, userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `
		UPDATE users
		SET last_connection = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
