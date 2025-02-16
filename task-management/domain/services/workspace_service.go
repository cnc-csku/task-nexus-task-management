package services

import (
	"context"
	"fmt"
	"math"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceService interface {
	SetupWorkspace(ctx context.Context, req *requests.CreateWorkspaceRequest, userID string) (*models.Workspace, *errutils.Error)
	ListOwnWorkspace(ctx context.Context, userId string) ([]responses.ListOwnWorkspaceResponseWorkspace, *errutils.Error)
	ListWorkspaceMembers(ctx context.Context, req *requests.ListWorkspaceMemberRequest) (*responses.ListWorkspaceMembersResponse, *errutils.Error)
}

type workspaceServiceImpl struct {
	workspaceRepo       repositories.WorkspaceRepository
	globalSettingRepo   repositories.GlobalSettingRepository
	userRepo            repositories.UserRepository
	workspaceMemberRepo repositories.WorkspaceMemberRepository
}

func NewWorkspaceService(
	workspaceRepo repositories.WorkspaceRepository,
	globalSettingRepo repositories.GlobalSettingRepository,
	userRepo repositories.UserRepository,
	workspaceMemberRepo repositories.WorkspaceMemberRepository,
) WorkspaceService {
	return &workspaceServiceImpl{
		workspaceRepo:       workspaceRepo,
		globalSettingRepo:   globalSettingRepo,
		userRepo:            userRepo,
		workspaceMemberRepo: workspaceMemberRepo,
	}
}

func (w *workspaceServiceImpl) SetupWorkspace(ctx context.Context, req *requests.CreateWorkspaceRequest, userID string) (*models.Workspace, *errutils.Error) {
	// Check is setup workspace
	isSetupWorkspace, err := w.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupWorkspace)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	if isSetupWorkspace == nil {
		err := w.globalSettingRepo.Set(ctx, &models.GlobalSetting{
			Key:   constant.GlobalSettingKeyIsSetupWorkspace,
			Type:  models.GlobalSettingTypeBool,
			Value: false,
		})

		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
	} else if isSetupWorkspace.Value.(bool) {
		return nil, errutils.NewError(exceptions.ErrWorkspaceAlreadySetup, errutils.BadRequest).
			WithMessage("Workspace already setup")
	}

	// Check is setup owner
	isSetupOwner, err := w.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupOwner)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	if isSetupOwner == nil {
		err := w.globalSettingRepo.Set(ctx, &models.GlobalSetting{
			Key:   constant.GlobalSettingKeyIsSetupOwner,
			Type:  models.GlobalSettingTypeBool,
			Value: false,
		})

		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
	}

	if !isSetupOwner.Value.(bool) {
		return nil, errutils.NewError(exceptions.ErrOwnerNotSetup, errutils.BadRequest)
	}

	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	// Find user
	user, err := w.userRepo.FindByID(ctx, userObjID)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	if user == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.BadRequest)
	}

	// Create workspace with user as owner
	workspace, err := w.workspaceRepo.Create(ctx, &repositories.CreateWorkspaceRequest{
		Name:            req.Name,
		UserDisplayName: user.DisplayName,
		ProfileUrl:      user.ProfileUrl,
	})

	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	_, err = w.workspaceMemberRepo.Create(ctx, &repositories.CreateWorkspaceMemberRequest{
		WorkspaceID: workspace.ID,
		UserID:      userObjID,
		Role:        models.WorkspaceMemberRoleOwner,
	})

	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	// Set is setup complete
	err = w.globalSettingRepo.Set(ctx, &models.GlobalSetting{
		Key:   constant.GlobalSettingKeyIsSetupWorkspace,
		Type:  models.GlobalSettingTypeBool,
		Value: true,
	})
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return workspace, nil
}

func (s *workspaceServiceImpl) ListOwnWorkspace(ctx context.Context, userId string) ([]responses.ListOwnWorkspaceResponseWorkspace, *errutils.Error) {
	userObjID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	workspaceMembers, err := s.workspaceMemberRepo.FindByUserID(ctx, userObjID)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if workspaceMembers == nil {
		return []responses.ListOwnWorkspaceResponseWorkspace{}, nil
	}

	workspaceIDs := make([]bson.ObjectID, 0)
	for _, workspaceMember := range workspaceMembers {
		workspaceIDs = append(workspaceIDs, workspaceMember.WorkspaceID)
	}

	workspaces, err := s.workspaceRepo.FindByWorkspaceIDs(ctx, workspaceIDs)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	workspaceMap := make(map[bson.ObjectID]models.Workspace)
	for _, workspace := range workspaces {
		workspaceMap[workspace.ID] = workspace
	}

	workspaceResponses := make([]responses.ListOwnWorkspaceResponseWorkspace, 0)
	for _, workspaceMember := range workspaceMembers {
		if workspace, exists := workspaceMap[workspaceMember.WorkspaceID]; exists {
			workspaceResponses = append(workspaceResponses, responses.ListOwnWorkspaceResponseWorkspace{
				ID:       workspace.ID.Hex(),
				Name:     workspace.Name,
				Role:     workspaceMember.Role.String(),
				JoinedAt: workspaceMember.JoinedAt,
			})
		}
	}

	return workspaceResponses, nil
}

func validateListWorkspaceMembersPaginationRequestSortBy(sortBy string) bool {
	switch sortBy {
	case constant.UserFieldFullName, constant.UserFieldDisplayName, constant.UserFieldEmail:
		return true
	}
	return false
}

func normalizeListWorkspaceMembersPaginationParams(req *requests.ListWorkspaceMemberRequest) {
	if req.PaginationRequest.Page <= 0 {
		req.PaginationRequest.Page = 1
	}
	if req.PaginationRequest.PageSize <= 0 {
		req.PaginationRequest.PageSize = 100
	}
	if req.PaginationRequest.SortBy == "" || !validateListForWorkspaceOwnerPaginationRequestSortBy(req.PaginationRequest.SortBy) {
		req.PaginationRequest.SortBy = constant.UserFieldDisplayName
	}
	if req.PaginationRequest.Order == "" {
		req.PaginationRequest.Order = constant.DESC
	}
}

func (s *workspaceServiceImpl) ListWorkspaceMembers(ctx context.Context, req *requests.ListWorkspaceMemberRequest) (*responses.ListWorkspaceMembersResponse, *errutils.Error) {
	fmt.Println(req.WorkspaceID)
	workspaceObjID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	normalizeListWorkspaceMembersPaginationParams(req)

	workspaceMembers, err := s.workspaceMemberRepo.FindByWorkspaceID(ctx, workspaceObjID)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	userIDs := make([]bson.ObjectID, 0)
	for _, member := range workspaceMembers {
		userIDs = append(userIDs, member.UserID)
	}

	users, totalUser, err := s.userRepo.SearchWithUserIDs(ctx, &repositories.SearchUserWithUserIDsRequest{
		UserIDs: userIDs,
		Keyword: req.Keyword,
		PaginationRequest: repositories.PaginationRequest{
			Page:     req.Page,
			PageSize: req.PageSize,
			SortBy:   req.SortBy,
			Order:    req.Order,
		},
	})
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	userMap := make(map[bson.ObjectID]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	response := make([]responses.ListWorkspaceMembersResponseWorkspaceMember, 0)
	for _, member := range workspaceMembers {
		if user, exists := userMap[member.UserID]; exists {
			response = append(response, responses.ListWorkspaceMembersResponseWorkspaceMember{
				WorkspaceMemberID: member.ID.Hex(),
				UserID:            user.ID.Hex(),
				Role:              member.Role.String(),
				JoinedAt:          member.JoinedAt,
				Email:             user.Email,
				FullName:          user.FullName,
				DisplayName:       user.DisplayName,
				ProfileUrl:        user.ProfileUrl,
			})
		}
	}

	return &responses.ListWorkspaceMembersResponse{
		Members: response,
		PaginationResponse: responses.PaginationResponse{
			Page:      req.Page,
			PageSize:  req.PageSize,
			TotalPage: int(math.Ceil(float64(totalUser) / float64(req.PaginationRequest.PageSize))),
			TotalItem: int(totalUser),
		},
	}, nil
}
