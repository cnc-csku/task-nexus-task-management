package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskCommentRepository interface {
	Create(ctx context.Context, taskComment *CreateTaskCommentRequest) (*models.TaskComment, error)
	FindByTaskID(ctx context.Context, taskID string) ([]*models.TaskComment, error)
}

type CreateTaskCommentRequest struct {
	TaskID  string
	Content string
	UserID  bson.ObjectID
}
