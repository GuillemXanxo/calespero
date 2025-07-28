package ports

import (
	"context"

	"calespero/internal/core/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) error
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	ValidateToken(token string) (string, error)
}
