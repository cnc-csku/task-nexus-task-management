package repositories

import (
	"context"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SprintRepository interface {
	Create(ctx context.Context, sprint *CreateSprintRequest) (*models.Sprint, error)
	FindByID(ctx context.Context, sprintID bson.ObjectID) (*models.Sprint, error)
	Update(ctx context.Context, sprint *UpdateSprintRequest) error
	FindByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.Sprint, error)
	List(ctx context.Context, filter *ListSprintFilter) ([]models.Sprint, error)
	UpdateStatus(ctx context.Context, req *UpdateSprintStatusRequest) (*models.Sprint, error)
	Delete(ctx context.Context, sprintID bson.ObjectID) error
}

type CreateSprintRequest struct {
	ProjectID bson.ObjectID
	Title     string
	CreatedBy bson.ObjectID
}

type UpdateSprintRequest struct {
	ID         bson.ObjectID
	Title      string
	SprintGoal string
	StartDate  *time.Time
	EndDate    *time.Time
	UpdatedBy  bson.ObjectID
}

type ListSprintFilter struct {
	ProjectID bson.ObjectID
	IsActive  *bool
}

type UpdateSprintStatusRequest struct {
	ID        bson.ObjectID
	Status    models.SprintStatus
	UpdatedBy bson.ObjectID
}
