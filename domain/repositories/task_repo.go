package repositories

import (
	"context"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskRepository interface {
	Create(ctx context.Context, task *CreateTaskRequest) (*models.Task, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Task, error)
	FindByIDs(ctx context.Context, ids []bson.ObjectID) ([]*models.Task, error)
	FindByTaskRef(ctx context.Context, taskRef string) (*models.Task, error)
	FindByTaskRefAndProjectID(ctx context.Context, taskRef string, projectID bson.ObjectID) (*models.Task, error)
	UpdateDetail(ctx context.Context, in *UpdateTaskDetailRequest) (*models.Task, error)
	UpdateTitle(ctx context.Context, in *UpdateTaskTitleRequest) (*models.Task, error)
	UpdateParentID(ctx context.Context, in *UpdateTaskParentIDRequest) (*models.Task, error)
	UpdateStatus(ctx context.Context, in *UpdateTaskStatusRequest) (*models.Task, error)
	UpdateApprovals(ctx context.Context, in *UpdateTaskApprovalsRequest) (*models.Task, error)
	ApproveTask(ctx context.Context, in *ApproveTaskRequest) (*models.Task, error)
	UpdateAssignees(ctx context.Context, in *UpdateTaskAssigneesRequest) (*models.Task, error)
	UpdateCurrentSprintID(ctx context.Context, in *UpdateTaskCurrentSprintIDRequest) (*models.Task, error)
	UpdatePreviousSprintIDs(ctx context.Context, in *UpdateTaskPreviousSprintIDsRequest) (*models.Task, error)
	FindByProjectIDAndStatuses(ctx context.Context, projectID bson.ObjectID, statuses []string) ([]*models.Task, error)
	UpdateHasChildren(ctx context.Context, in *UpdateTaskHasChildrenRequest) (*models.Task, error)
	FindByParentID(ctx context.Context, parentID bson.ObjectID) ([]*models.Task, error)
	UpdateChildrenPoint(ctx context.Context, in *UpdateTaskChildrenPointRequest) (*models.Task, error)
	FindByProjectIDAndType(ctx context.Context, projectID bson.ObjectID, taskType models.TaskType) ([]*models.Task, error)
	Search(ctx context.Context, in *SearchTaskRequest) ([]*models.Task, error)
	UpdateAttributes(ctx context.Context, in *UpdateTaskAttributesRequest) (*models.Task, error)
	FindByCurrentSprintID(ctx context.Context, sprintID bson.ObjectID) ([]*models.Task, error)
	FindByPreviousSprintID(ctx context.Context, sprintID bson.ObjectID) ([]*models.Task, error)
	FindByCurrentSprintIDAndPreviousSprintIDs(ctx context.Context, sprintID bson.ObjectID) ([]*models.Task, error)
}

type CreateTaskRequest struct {
	TaskRef     string
	ProjectID   bson.ObjectID
	Title       string
	Description string
	ParentID    *bson.ObjectID
	Type        models.TaskType
	Status      string
	Priority    models.TaskPriority
	Sprint      *models.TaskSprint
	StartDate   *time.Time
	DueDate     *time.Time
	CreatedBy   bson.ObjectID
}

type UpdateTaskDetailRequest struct {
	ID          bson.ObjectID
	Title       string
	Description string
	Priority    string
	StartDate   *time.Time
	DueDate     *time.Time
	UpdatedBy   bson.ObjectID
}

type UpdateTaskTitleRequest struct {
	ID        bson.ObjectID
	Title     string
	UpdatedBy bson.ObjectID
}

type UpdateTaskParentIDRequest struct {
	ID        bson.ObjectID
	ParentID  *bson.ObjectID
	UpdatedBy bson.ObjectID
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
	Point    *int
}

type UpdateTaskCurrentSprintIDRequest struct {
	ID              bson.ObjectID
	CurrentSprintID *bson.ObjectID
	UpdatedBy       bson.ObjectID
}

type UpdateTaskPreviousSprintIDsRequest struct {
	ID                bson.ObjectID
	PreviousSprintIDs []bson.ObjectID
	UpdatedBy         bson.ObjectID
}

type UpdateTaskHasChildrenRequest struct {
	ID          bson.ObjectID
	HasChildren bool
}

type UpdateTaskChildrenPointRequest struct {
	ID            bson.ObjectID
	ChildrenPoint int
}

type SearchTaskRequest struct {
	ProjectID          bson.ObjectID
	TaskTypes          []models.TaskType
	SprintID           *bson.ObjectID
	IsTaskWithNoSprint bool
	EpicTaskID         *bson.ObjectID
	IsTaskWithNoEpic   bool
	UserIDs            []bson.ObjectID
	Positions          []string
	Statuses           []string
	IsDoneStatuses     []string
	SearchKeyword      *string
}

type UpdateTaskAttributesRequest struct {
	ID         bson.ObjectID
	Attributes []models.TaskAttribute
	UpdatedBy  bson.ObjectID
}
