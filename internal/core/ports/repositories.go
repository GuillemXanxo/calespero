package ports

import (
	"context"

	"calespero/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateLastConnection(ctx context.Context, userID string) error
}
