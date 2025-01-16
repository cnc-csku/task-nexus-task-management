package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InvitationService interface {
	Create(ctx context.Context, req *requests.CreateInvitationRequest, inviterUserID string) (*responses.CreateInvitationResponse, *errutils.Error)
	ListForUser(ctx context.Context, userID string) (*responses.ListInvitationForUserResponse, *errutils.Error)
	UserResponse(ctx context.Context, req *requests.UserResponseInvitationRequest, userID string) (*responses.UserResponseInvitationResponse, *errutils.Error)
}

type invitationServiceImpl struct {
	userRepo       repositories.UserRepository
	workspaceRepo  repositories.WorkspaceRepository
	invitationRepo repositories.InvitationRepository
	config         *config.Config
}

func NewInvitationService(userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository, invitationRepo repositories.InvitationRepository, config *config.Config) InvitationService {
	return &invitationServiceImpl{
		userRepo:       userRepo,
		workspaceRepo:  workspaceRepo,
		invitationRepo: invitationRepo,
		config:         config,
	}
}

func (i *invitationServiceImpl) Create(ctx context.Context, req *requests.CreateInvitationRequest, inviterUserID string) (*responses.CreateInvitationResponse, *errutils.Error) {
	bsonWorkspaceID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonInviterUserID, err := bson.ObjectIDFromHex(inviterUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonInviteeUserID, err := bson.ObjectIDFromHex(req.InviteeUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the inviter is the admin of the workspace
	inviter, err := i.workspaceRepo.FindWorkspaceMemberByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonInviterUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if inviter == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInWorkspace, errutils.BadRequest).WithDebugMessage("Inviter not found in workspace")
	} else if inviter.Role != models.WorkspaceMemberRoleAdmin {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("Inviter is not an admin")
	}

	// Check if the invitee is already a member of the workspace
	invitee, err := i.workspaceRepo.FindWorkspaceMemberByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonInviteeUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if invitee != nil {
		return nil, errutils.NewError(exceptions.ErrMemberAlreadyInWorkspace, errutils.BadRequest).WithDebugMessage("Invitee is already a member of the workspace")
	}

	// Check if the invitee is already invited to the workspace
	invitation, err := i.invitationRepo.FindByWorkspaceIDAndInviteeUserID(ctx, bsonWorkspaceID, bsonInviteeUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if invitation != nil {
		return nil, errutils.NewError(exceptions.ErrInvitationAlreadySent, errutils.BadRequest).WithDebugMessage("Invitee is already invited to the workspace")
	}

	// Create the invitation
	createInvitationReq := &repositories.CreateInvitationRequest{
		WorkspaceID:   bsonWorkspaceID,
		InviteeUserID: bsonInviteeUserID,
		Status:        models.InvitationStatusPending,
		ExpiredAt:     time.Now().Add(constant.InvitationExpirationIn),
		CustomMessage: req.CustomMessage,
		CreatedBy:     bsonInviterUserID,
	}

	err = i.invitationRepo.Create(ctx, createInvitationReq)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return &responses.CreateInvitationResponse{
		Message: "Invitation sent successfully",
	}, nil
}

func (i *invitationServiceImpl) ListForUser(ctx context.Context, userID string) (*responses.ListInvitationForUserResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	invitations, err := i.invitationRepo.FindByInviteeUserID(ctx, bsonUserID, constant.InvitationFieldCreatedAt, constant.DESC)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}
	for i, invitation := range invitations {
		fmt.Println("\n", i, invitation)
	}

	invitationResponses := make([]responses.InvitationForUserResponse, 0)
	for _, invitation := range invitations {
		workspace, err := i.workspaceRepo.FindByID(ctx, invitation.WorkspaceID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		inviter, err := i.userRepo.FindByID(ctx, invitation.CreatedBy)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		var status = invitation.Status
		if invitation.ExpiredAt.Before(time.Now()) {
			status = models.InvitationStatusExpired
		}

		invitationResponses = append(invitationResponses, responses.InvitationForUserResponse{
			InvitationID:       invitation.ID.Hex(),
			WorkspaceID:        invitation.WorkspaceID.Hex(),
			WorkspaceName:      workspace.Name,
			Status:             status.String(),
			CustomMessage:      invitation.CustomMessage,
			InvitedAt:          invitation.CreatedAt.Format(constant.TimeFormat),
			InviterDisplayName: inviter.DisplayName,
			InviterFullName:    inviter.FullName,
			InviterUserID:      invitation.CreatedBy.Hex(),
			ExpiredAt:          invitation.ExpiredAt.Format(constant.TimeFormat),
			IsExpired:          time.Now().After(invitation.ExpiredAt),
			RespondedAt:        invitation.RespondedAt,
		})
	}

	return &responses.ListInvitationForUserResponse{
		Invitations: invitationResponses,
	}, nil
}

func (i *invitationServiceImpl) UserResponse(ctx context.Context, req *requests.UserResponseInvitationRequest, userID string) (*responses.UserResponseInvitationResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonInvitationID, err := bson.ObjectIDFromHex(req.InvitationID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	invitation, err := i.invitationRepo.FindByID(ctx, bsonInvitationID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if invitation == nil {
		return nil, errutils.NewError(exceptions.ErrInvitationNotFound, errutils.BadRequest).WithDebugMessage("Invitation not found")
	}

	if invitation.InviteeUserID.Hex() != bsonUserID.Hex() {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("Permission denied")
	}

	if invitation.Status != models.InvitationStatusPending {
		return nil, errutils.NewError(exceptions.ErrInvalidInvitationStatus, errutils.BadRequest).WithDebugMessage("Invalid invitation status")
	}

	user, err := i.userRepo.FindByID(ctx, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	if req.Action == constant.InvitationActionAccept {
		// Add the invitee as a member of the workspace
		createWorkspaceMemberReq := &repositories.CreateWorkspaceMemberRequest{
			WorkspaceID: invitation.WorkspaceID,
			UserID:      invitation.InviteeUserID,
			Name:        user.FullName,
			Role:        models.WorkspaceMemberRoleUser,
		}

		err = i.workspaceRepo.CreateWorkspaceMember(ctx, createWorkspaceMemberReq)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		// Update the invitation status to accepted
		err = i.invitationRepo.UpdateStatus(ctx, bsonInvitationID, models.InvitationStatusAccepted)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		return &responses.UserResponseInvitationResponse{
			Message: "Invitation accepted successfully",
		}, nil
	} else if req.Action == constant.InvitationActionDecline {
		// Update the invitation status to declined
		err = i.invitationRepo.UpdateStatus(ctx, bsonInvitationID, models.InvitationStatusDeclined)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		return &responses.UserResponseInvitationResponse{
			Message: "Invitation declined successfully",
		}, nil
	} else {
		return nil, errutils.NewError(exceptions.ErrInvalidInvitationAction, errutils.BadRequest).WithDebugMessage("Invalid invitation action")
	}
}
