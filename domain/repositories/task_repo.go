package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskRepository interface {
	Create(ctx context.Context, task *CreateTaskRequest) (*models.Task, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Task, error)
	FindByTaskRef(ctx context.Context, taskRef string) (*models.Task, error)
	UpdateDetail(ctx context.Context, in *UpdateTaskDetailRequest) (*models.Task, error)
	UpdateStatus(ctx context.Context, in *UpdateTaskStatusRequest) (*models.Task, error)
	UpdateApprovals(ctx context.Context, in *UpdateTaskApprovalsRequest) (*models.Task, error)
	ApproveTask(ctx context.Context, in *ApproveTaskRequest) (*models.Task, error)
	UpdateAssignees(ctx context.Context, in *UpdateTaskAssigneesRequest) (*models.Task, error)
	UpdateSprint(ctx context.Context, in *UpdateTaskSprintRequest) (*models.Task, error)
}

type CreateTaskRequest struct {
	TaskRef     string
	ProjectID   bson.ObjectID
	Title       string
	Description string
	ParentID    *bson.ObjectID
	Type        models.TaskType
	Status      string
	Sprint      *models.TaskSprint
	CreatedBy   bson.ObjectID
}

type UpdateTaskDetailRequest struct {
	ID          bson.ObjectID
	Title       string
	Description string
	ParentID    *bson.ObjectID
	Type        models.TaskType
	Priority    *string
	UpdatedBy   bson.ObjectID
}

type UpdateTaskStatusRequest struct {
	ID        bson.ObjectID
	Status    string
	UpdatedBy bson.ObjectID
}

type UpdateTaskApprovalsRequest struct {
	ID        bson.ObjectID
	Approval  []UpdateTaskApprovalsRequestApproval
	UpdatedBy bson.ObjectID
}

type UpdateTaskApprovalsRequestApproval struct {
	UserID bson.ObjectID
}

type ApproveTaskRequest struct {
	ID     bson.ObjectID
	Reason string
	UserID bson.ObjectID
}

type UpdateTaskAssigneesRequest struct {
	ID        bson.ObjectID
	Assignees []UpdateTaskAssigneesRequestAssignee
	UpdatedBy bson.ObjectID
}

type UpdateTaskAssigneesRequestAssignee struct {
	Position string
	UserID   bson.ObjectID
	Point    int
}

type UpdateTaskSprintRequest struct {
	ID                bson.ObjectID
	CurrentSprintID   bson.ObjectID
	PreviousSprintIDs []bson.ObjectID
	UpdatedBy         bson.ObjectID
}
