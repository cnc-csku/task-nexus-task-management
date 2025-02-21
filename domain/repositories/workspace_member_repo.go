package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceMemberRepository interface {
	FindByWorkspaceID(ctx context.Context, workspaceID bson.ObjectID) ([]models.WorkspaceMember, error)
	Create(ctx context.Context, req *CreateWorkspaceMemberRequest) (*models.WorkspaceMember, error)
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]models.WorkspaceMember, error)
	FindByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) (*models.WorkspaceMember, error)
}

type CreateWorkspaceMemberRequest struct {
	WorkspaceID bson.ObjectID
	UserID      bson.ObjectID
	Role        models.WorkspaceMemberRole
}
