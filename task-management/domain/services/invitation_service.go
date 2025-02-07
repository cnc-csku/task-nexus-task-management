package services

import (
	"context"
	"math"
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
	ListForWorkspaceOwner(ctx context.Context, req *requests.ListInvitationForWorkspaceOwnerParams, userID string) (*responses.ListInvitationForWorkspaceOwnerResponse, *errutils.Error)
	UserResponse(ctx context.Context, req *requests.UserResponseInvitationRequest, userID string) (*responses.UserResponseInvitationResponse, *errutils.Error)
}

type invitationServiceImpl struct {
	userRepo            repositories.UserRepository
	workspaceRepo       repositories.WorkspaceRepository
	invitationRepo      repositories.InvitationRepository
	workspaceMemberRepo repositories.WorkspaceMemberRepository
	config              *config.Config
}

func NewInvitationService(
	userRepo repositories.UserRepository,
	workspaceRepo repositories.WorkspaceRepository,
	invitationRepo repositories.InvitationRepository,
	workspaceMemberRepo repositories.WorkspaceMemberRepository,
	config *config.Config,
) InvitationService {
	return &invitationServiceImpl{
		userRepo:            userRepo,
		workspaceRepo:       workspaceRepo,
		invitationRepo:      invitationRepo,
		workspaceMemberRepo: workspaceMemberRepo,
		config:              config,
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

	// Check if the inviter is the owner of the workspace
	inviter, err := i.workspaceMemberRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonInviterUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if inviter == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInWorkspace, errutils.BadRequest).WithDebugMessage("Inviter not found in workspace")
	} else if inviter.Role != models.WorkspaceMemberRoleOwner {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("Inviter is not an owner")
	}

	// Check if the invitee is already a member of the workspace
	invitee, err := i.workspaceMemberRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonInviteeUserID)
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
		Role:          models.InvitationRole(req.Role),
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
		if status == models.InvitationStatusPending && invitation.ExpiredAt.Before(time.Now()) {
			status = models.InvitationStatusExpired
		}

		invitationResponses = append(invitationResponses, responses.InvitationForUserResponse{
			InvitationID:       invitation.ID.Hex(),
			WorkspaceID:        invitation.WorkspaceID.Hex(),
			WorkspaceName:      workspace.Name,
			Role:               invitation.Role.String(),
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

func validateListForWorkspaceOwnerPaginationRequestSortBy(sortBy string) bool {
	switch sortBy {
	case constant.InvitationFieldCreatedAt, constant.InvitationFieldStatus:
		return true
	}
	return false
}

func validateListForWorkspaceOwnerSearchBy(searchBy string) bool {
	switch searchBy {
	case constant.InvitationFieldStatus:
		return true
	}
	return false
}

func normalizeListForWorkspaceOwnerPaginationParams(req *requests.ListInvitationForWorkspaceOwnerParams) {
	if req.PaginationRequest.Page <= 0 {
		req.PaginationRequest.Page = 1
	}
	if req.PaginationRequest.PageSize <= 0 {
		req.PaginationRequest.PageSize = 100
	}
	if req.PaginationRequest.SortBy == "" || !validateListForWorkspaceOwnerPaginationRequestSortBy(req.PaginationRequest.SortBy) {
		req.PaginationRequest.SortBy = constant.InvitationFieldCreatedAt
	}
	if req.PaginationRequest.Order == "" {
		req.PaginationRequest.Order = constant.DESC
	}
}

func (i *invitationServiceImpl) ListForWorkspaceOwner(ctx context.Context, req *requests.ListInvitationForWorkspaceOwnerParams, userID string) (*responses.ListInvitationForWorkspaceOwnerResponse, *errutils.Error) {
	if req.SearchBy != "" || !validateListForWorkspaceOwnerSearchBy(req.SearchBy) {
		req.SearchBy = ""
	}

	normalizeListForWorkspaceOwnerPaginationParams(req)

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonWorkspaceID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the user is the owner of the workspace
	member, err := i.workspaceMemberRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInWorkspace, errutils.BadRequest).WithDebugMessage("User not found in workspace")
	} else if member.Role != models.WorkspaceMemberRoleOwner {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not an owner")
	}

	invitations, totalInvitation, err := i.invitationRepo.SearchInvitationForEachWorkspace(ctx, &repositories.SearchInvitationForEachWorkspaceRequest{
		WorkspaceID: bsonWorkspaceID,
		Keyword:     req.Keyword,
		SearchBy:    req.SearchBy,
		PaginationRequest: repositories.PaginationRequest{
			Page:     req.PaginationRequest.Page,
			PageSize: req.PaginationRequest.PageSize,
			SortBy:   req.PaginationRequest.SortBy,
			Order:    req.PaginationRequest.Order,
		},
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	invitationResponses := make([]responses.InvitationForWorkspaceOwnerResponse, 0)
	for _, invitation := range invitations {
		workspace, err := i.workspaceRepo.FindByID(ctx, invitation.WorkspaceID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		inviter, err := i.userRepo.FindByID(ctx, invitation.CreatedBy)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		invitee, err := i.userRepo.FindByID(ctx, invitation.InviteeUserID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		var status = invitation.Status
		if invitation.ExpiredAt.Before(time.Now()) {
			status = models.InvitationStatusExpired
		}

		invitationResponses = append(invitationResponses, responses.InvitationForWorkspaceOwnerResponse{
			InvitationID:       invitation.ID.Hex(),
			WorkspaceID:        invitation.WorkspaceID.Hex(),
			WorkspaceName:      workspace.Name,
			Role:               invitation.Role.String(),
			Status:             status.String(),
			CustomMessage:      invitation.CustomMessage,
			InvitedAt:          invitation.CreatedAt.Format(constant.TimeFormat),
			InviteeDisplayName: invitee.DisplayName,
			InviteeFullName:    invitee.FullName,
			InviteeUserID:      invitation.InviteeUserID.Hex(),
			InviterDisplayName: inviter.DisplayName,
			InviterFullName:    inviter.FullName,
			InviterUserID:      invitation.CreatedBy.Hex(),
			ExpiredAt:          invitation.ExpiredAt.Format(constant.TimeFormat),
			IsExpired:          time.Now().After(invitation.ExpiredAt),
			RespondedAt:        invitation.RespondedAt,
		})
	}

	return &responses.ListInvitationForWorkspaceOwnerResponse{
		Invitations: invitationResponses,
		PaginationResponse: responses.PaginationResponse{
			Page:      req.PaginationRequest.Page,
			PageSize:  req.PaginationRequest.PageSize,
			TotalPage: int(math.Ceil(float64(totalInvitation) / float64(req.PaginationRequest.PageSize))),
			TotalItem: int(totalInvitation),
		},
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
		return nil, errutils.NewError(exceptions.ErrInvitationAlreadyResponded, errutils.BadRequest).WithDebugMessage("Invitation already responded")
	}

	if req.Action == constant.InvitationActionAccept {
		var role models.WorkspaceMemberRole
		if invitation.Role == models.InvitationRoleModerator {
			role = models.WorkspaceMemberRoleModerator
		} else {
			role = models.WorkspaceMemberRoleMember
		}

		// Add the invitee as a member of the workspace
		createWorkspaceMemberReq := &repositories.CreateWorkspaceMemberRequest{
			WorkspaceID: invitation.WorkspaceID,
			UserID:      invitation.InviteeUserID,
			Role:        role,
		}

		_, err = i.workspaceMemberRepo.Create(ctx, createWorkspaceMemberReq)
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
