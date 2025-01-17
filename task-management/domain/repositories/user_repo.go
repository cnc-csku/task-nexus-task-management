package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *CreateUserRequest) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByIDs(ctx context.Context, userIDs []bson.ObjectID) ([]models.User, error)
	Search(ctx context.Context, in *SearchUserRequest) ([]*models.User, int64, error)
	FindByID(ctx context.Context, userID bson.ObjectID) (*models.User, error)
}

type CreateUserRequest struct {
	Email        string
	PasswordHash string
	FullName     string
	DisplayName  string
	ProfileUrl   string
}

type SearchUserRequest struct {
	Keyword           string
	PaginationRequest PaginationRequest
}
