package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskCommentRepository interface {
	Create(ctx context.Context, taskComment *CreateTaskCommentRequest) (*models.TaskComment, error)
	FindByTaskID(ctx context.Context, taskID bson.ObjectID) ([]*models.TaskComment, error)
}

type CreateTaskCommentRequest struct {
	TaskID  bson.ObjectID
	Content string
	UserID  bson.ObjectID
}
