package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepository interface {
	FindByWorkspaceIDAndName(ctx context.Context, workspaceID bson.ObjectID, name string) (*models.Project, error)
	FindByWorkspaceIDAndProjectPrefix(ctx context.Context, workspaceID bson.ObjectID, projectPrefix string) (*models.Project, error)
	Create(ctx context.Context, project *CreateProjectRequest) (*models.Project, error)
}

type CreateProjectRequest struct {
	WorkspaceID   bson.ObjectID
	Name          string
	ProjectPrefix string
	Description   string
	Status        string
	Members       []models.Member
	CreatedBy     bson.ObjectID
}
