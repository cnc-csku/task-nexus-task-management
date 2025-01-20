package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepository interface {
	FindByProjectID(ctx context.Context, projectID bson.ObjectID) (*models.Project, error)
	FindByWorkspaceIDAndName(ctx context.Context, workspaceID bson.ObjectID, name string) (*models.Project, error)
	FindByWorkspaceIDAndProjectPrefix(ctx context.Context, workspaceID bson.ObjectID, projectPrefix string) (*models.Project, error)
	Create(ctx context.Context, project *CreateProjectRequest) (*models.Project, error)
	FindByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) ([]*models.Project, error)
	FindMemberByProjectIDAndUserID(ctx context.Context, projectID bson.ObjectID, userID bson.ObjectID) (*models.ProjectMember, error)
	FindPositionByProjectID(ctx context.Context, projectID bson.ObjectID) ([]string, error)
	AddPosition(ctx context.Context, projectID bson.ObjectID, position []string) error
}

type CreateProjectRequest struct {
	WorkspaceID   bson.ObjectID
	Name          string
	ProjectPrefix string
	Description   string
	Status        models.ProjectStatus
	Owner         *models.ProjectMember
	// Members       []models.ProjectMember
	CreatedBy bson.ObjectID
}
