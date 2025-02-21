package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceRepository interface {
	FindByID(ctx context.Context, workspaceID bson.ObjectID) (*models.Workspace, error)
	Create(ctx context.Context, workspace *CreateWorkspaceRequest) (*models.Workspace, error)
	FindByWorkspaceIDs(ctx context.Context, workspaceIDs []bson.ObjectID) ([]models.Workspace, error)
}

type CreateWorkspaceRequest struct {
	UserID          bson.ObjectID
	Name            string
	UserDisplayName string
	ProfileUrl      string
}
