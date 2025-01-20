package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceRepository interface {
	FindWorkspaceMemberByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) (*models.WorkspaceMember, error)
	FindByID(ctx context.Context, workspaceID bson.ObjectID) (*models.Workspace, error)
	CreateWorkspaceMember(ctx context.Context, req *CreateWorkspaceMemberRequest) error
	Create(ctx context.Context, workspace *CreateWorkspaceRequest) (*models.Workspace, error)
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]*models.Workspace, error)
}

type CreateWorkspaceMemberRequest struct {
	WorkspaceID bson.ObjectID
	UserID      bson.ObjectID
	Name        string
	Role        models.WorkspaceMemberRole
}

type CreateWorkspaceRequest struct {
	Name     string
	UserID   bson.ObjectID
	UserName string
}
