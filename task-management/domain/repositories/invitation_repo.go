package repositories

import (
	"context"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InvitationRepository interface {
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Invitation, error)
	FindByWorkspaceIDAndInviteeUserID(ctx context.Context, workspaceID bson.ObjectID, inviteeUserID bson.ObjectID) (*models.Invitation, error)
	Create(ctx context.Context, invitation *CreateInvitationRequest) error
	FindByInviteeUserID(ctx context.Context, inviteeUserID bson.ObjectID, sortBy string, order string) ([]models.Invitation, error)
	UpdateStatus(ctx context.Context, id bson.ObjectID, status models.InvitationStatus) error
	SearchInvitationForEachWorkspaceRequest(ctx context.Context, in *SearchInvitationForEachWorkspaceRequest) ([]models.Invitation, int64, error)
}

type CreateInvitationRequest struct {
	WorkspaceID   bson.ObjectID
	InviteeUserID bson.ObjectID
	Role          models.InvitationRole
	Status        models.InvitationStatus
	ExpiredAt     time.Time
	CustomMessage string
	CreatedBy     bson.ObjectID
}

type SearchInvitationForEachWorkspaceRequest struct {
	WorkspaceID       bson.ObjectID
	Keyword           string
	SearchBy          string
	PaginationRequest PaginationRequest
}
