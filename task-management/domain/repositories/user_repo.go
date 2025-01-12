package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *CreateUserRequest) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type CreateUserRequest struct {
	Email        string 
	PasswordHash string
	FullName     string
	DisplayName  string
}
