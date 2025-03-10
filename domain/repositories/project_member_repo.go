package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectMemberRepository interface {
	Create(ctx context.Context, projectMember *CreateProjectMemberRequest) error
	CreateMany(ctx context.Context, projectMembers []CreateProjectMemberRequest) error
	FindByID(ctx context.Context, id bson.ObjectID) (*models.ProjectMember, error)
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]*models.ProjectMember, error)
	FindByProjectID(ctx context.Context, projectID bson.ObjectID) ([]*models.ProjectMember, error)
	FindByProjectIDAndUserID(ctx context.Context, projectID bson.ObjectID, userID bson.ObjectID) (*models.ProjectMember, error)
	FindProjectOwnerByProjectID(ctx context.Context, projectID bson.ObjectID) (*models.ProjectMember, error)
	FindProjectOwnersByProjectIDs(ctx context.Context, projectIDs []bson.ObjectID) (map[bson.ObjectID]models.ProjectMember, error)
	FindByProjectIDAndPositions(ctx context.Context, projectID bson.ObjectID, positions []string) ([]*models.ProjectMember, error)
	UpdatePositionByID(ctx context.Context, in *UpdatePositionRequest) (*models.ProjectMember, error)
}

type CreateProjectMemberRequest struct {
	UserID    bson.ObjectID
	ProjectID bson.ObjectID
	Role      models.ProjectMemberRole
	Position  string
}

type SearchProjectMemberRequest struct {
	ProjectID         bson.ObjectID
	Keyword           string
	PaginationRequest PaginationRequest
}

type UpdatePositionRequest struct {
	ID       bson.ObjectID
	Position string
}
