package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskRepository interface {
	Create(ctx context.Context, task *CreateTaskRequest) (*models.Task, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Task, error)
	FindByTaskID(ctx context.Context, taskID string) (*models.Task, error)
}

type CreateTaskRequest struct {
	TaskID      string
	ProjectID   bson.ObjectID
	Title       string
	Description *string
	ParentID    *string
	Type        models.TaskType
	Status      string
	Sprint      *models.TaskSprint
	CreatedBy   bson.ObjectID
}
