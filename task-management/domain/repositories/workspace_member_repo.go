package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceMemberRepository interface {
	FindByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) (*models.WorkspaceMember, error)
	Create(ctx context.Context, req *CreateWorkspaceMemberRequest) error
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]models.WorkspaceMember, error)
	FindByWorkspaceID(ctx context.Context, workspaceID bson.ObjectID) ([]models.WorkspaceMember, error)
}

type CreateWorkspaceMemberRequest struct {
	UserID      bson.ObjectID
	WorkspaceID bson.ObjectID
	Role        models.WorkspaceMemberRole
}
