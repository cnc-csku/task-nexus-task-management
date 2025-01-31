package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RunningNumberRepository interface {
	Create(ctx context.Context, in *CreateRunningNumberRequest) (*models.RunningNumber, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.RunningNumber, error)
	IncrementSequence(ctx context.Context, id bson.ObjectID) error
}

type CreateRunningNumberRequest struct {
	Type models.RunningNumberType
}
