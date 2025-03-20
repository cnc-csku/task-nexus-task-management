package services

import (
	"context"
	"math"
	"strings"

	"github.com/cnc-csku/task-nexus-go-lib/utils/array"
	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectService interface {
	Create(ctx context.Context, req *requests.CreateProjectRequest, userID string) (*responses.CreateProjectResponse, *errutils.Error)
	ListMyProjects(ctx context.Context, req *requests.ListMyProjectsPathParams, userID string) ([]responses.ListProjectsResponse, *errutils.Error)
	GetProjectDetail(ctx context.Context, req *requests.GetProjectsDetailPathParams, userID string) (*responses.GetProjectDetailResponse, *errutils.Error)
	UpdateSetupStatus(ctx context.Context, req *requests.UpdateSetupStatusRequest, userID string) (*models.Project, *errutils.Error)
	UpdatePositions(ctx context.Context, req *requests.UpdatePositionsRequest, userID string) (*responses.UpdatePositionsResponse, *errutils.Error)
	ListPositions(ctx context.Context, req *requests.ListPositionsPathParams) ([]string, *errutils.Error)
	AddMembers(ctx context.Context, req *requests.AddProjectMembersRequest, userID string) (*responses.AddProjectMembersResponse, *errutils.Error)
	ListMembers(ctx context.Context, req *requests.ListProjectMembersRequest) (*responses.ListProjectMembersResponse, *errutils.Error)
	UpdateWorkflows(ctx context.Context, req *requests.UpdateWorkflowsRequest, userID string) (*responses.UpdateWorkflowsResponse, *errutils.Error)
	ListWorkflows(ctx context.Context, req *requests.ListWorkflowsPathParams) ([]models.ProjectWorkflow, *errutils.Error)
	UpdateAttributeTemplates(ctx context.Context, req *requests.UpdateAttributeTemplatesRequest, userID string) (*responses.UpdateAttributeTemplatesResponse, *errutils.Error)
	ListAttributeTemplates(ctx context.Context, req *requests.ListAttributeTemplatesPathParams) ([]models.ProjectAttributeTemplate, *errutils.Error)
}

type projectServiceImpl struct {
	userRepo            repositories.UserRepository
	workspaceRepo       repositories.WorkspaceRepository
	workspaceMemberRepo repositories.WorkspaceMemberRepository
	projectRepo         repositories.ProjectRepository
	projectMemberRepo   repositories.ProjectMemberRepository
	config              *config.Config
	taskRepo            repositories.TaskRepository
}

func NewProjectService(
	userRepo repositories.UserRepository,
	workspaceRepo repositories.WorkspaceRepository,
	workspaceMemberRepo repositories.WorkspaceMemberRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	config *config.Config,
	taskRepo repositories.TaskRepository,
) ProjectService {
	return &projectServiceImpl{
		userRepo:            userRepo,
		workspaceRepo:       workspaceRepo,
		workspaceMemberRepo: workspaceMemberRepo,
		projectRepo:         projectRepo,
		projectMemberRepo:   projectMemberRepo,
		config:              config,
		taskRepo:            taskRepo,
	}
}

func (p *projectServiceImpl) Create(ctx context.Context, req *requests.CreateProjectRequest, userID string) (*responses.CreateProjectResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonWorkspaceID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInvalidWorkspaceID, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the creator is owner or moderator of the workspace
	member, err := p.workspaceMemberRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInWorkspace, errutils.BadRequest)
	} else if member.Role != models.WorkspaceMemberRoleOwner && member.Role != models.WorkspaceMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	// Check if project's name already exists
	existsProjectByName, err := p.projectRepo.FindByWorkspaceIDAndName(ctx, bsonWorkspaceID, req.Name)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}
	if existsProjectByName != nil {
		return nil, errutils.NewError(exceptions.ErrProjectNameAlreadyExists, errutils.BadRequest)
	}

	// Check if project's prefix already exists
	existsProjectByPrefix, err := p.projectRepo.FindByWorkspaceIDAndProjectPrefix(ctx, bsonWorkspaceID, req.ProjectPrefix)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}
	if existsProjectByPrefix != nil {
		return nil, errutils.NewError(exceptions.ErrProjectPrefixAlreadyExists, errutils.BadRequest)
	}

	// Create owner member
	owner := &models.ProjectMember{
		UserID: member.UserID,
		Role:   models.ProjectMemberRoleOwner,
	}

	project := &repositories.CreateProjectRequest{
		WorkspaceID:        bsonWorkspaceID,
		Name:               req.Name,
		ProjectPrefix:      req.ProjectPrefix,
		Description:        req.Description,
		Status:             models.ProjectStatusActive,
		Owner:              owner,
		Workflows:          models.GetDefaultWorkflows(),
		Positions:          models.GetDefaultPositions(),
		AttributeTemplates: []models.ProjectAttributeTemplate{},
		CreatedBy:          bsonUserID,
	}

	createdProject, err := p.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	res := &responses.CreateProjectResponse{
		ID:            createdProject.ID.Hex(),
		WorkspaceID:   createdProject.WorkspaceID.Hex(),
		Name:          createdProject.Name,
		ProjectPrefix: createdProject.ProjectPrefix,
		Description:   createdProject.Description,
	}

	return res, nil
}

func (p *projectServiceImpl) ListMyProjects(ctx context.Context, req *requests.ListMyProjectsPathParams, userID string) ([]responses.ListProjectsResponse, *errutils.Error) {
	bsonWorkspaceID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInvalidWorkspaceID, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	projectMembers, err := p.projectMemberRepo.FindByUserID(ctx, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	if len(projectMembers) == 0 {
		return []responses.ListProjectsResponse{}, nil
	}

	bsonProjectIDs := []bson.ObjectID{}
	for _, projectMember := range projectMembers {
		bsonProjectIDs = append(bsonProjectIDs, projectMember.ProjectID)
	}

	projects, err := p.projectRepo.FindByProjectIDsAndWorkspaceID(ctx, bsonProjectIDs, bsonWorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	owners, err := p.projectMemberRepo.FindProjectOwnersByProjectIDs(ctx, bsonProjectIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	var bsonOwnerIDs []bson.ObjectID
	for _, owner := range owners {
		bsonOwnerIDs = append(bsonOwnerIDs, owner.UserID)
	}

	ownerInfo, err := p.userRepo.FindByIDs(ctx, bsonOwnerIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	ownerMap := make(map[bson.ObjectID]models.User)
	for _, owner := range ownerInfo {
		ownerMap[owner.ID] = owner
	}

	resp := make([]responses.ListProjectsResponse, 0)
	for _, project := range projects {
		owner, ok := owners[project.ID]
		if !ok {
			continue
		}
		ownerInfo, ok := ownerMap[owner.UserID]
		if !ok {
			continue
		}

		var profileUrl = ownerInfo.DefaultProfileUrl
		if ownerInfo.UploadedProfileUrl != nil {
			profileUrl = *ownerInfo.UploadedProfileUrl
		}

		resp = append(resp, responses.ListProjectsResponse{
			ID:                   project.ID.Hex(),
			WorkspaceID:          project.WorkspaceID.Hex(),
			Name:                 project.Name,
			ProjectPrefix:        project.ProjectPrefix,
			Description:          project.Description,
			Status:               project.Status.String(),
			OwnerUserID:          owner.UserID.Hex(),
			OwnerProjectMemberID: owner.ID.Hex(),
			OwnerDisplayName:     ownerInfo.DisplayName,
			OwnerProfileUrl:      profileUrl,
			CreatedAt:            project.CreatedAt,
			CreatedBy:            project.CreatedBy.Hex(),
			UpdatedAt:            project.UpdatedAt,
			UpdatedBy:            project.UpdatedBy.Hex(),
		})
	}

	return resp, nil
}

func (p *projectServiceImpl) GetProjectDetail(ctx context.Context, req *requests.GetProjectsDetailPathParams, userID string) (*responses.GetProjectDetailResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the user is a member of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound).WithDebugMessage("Project not found")
	}

	owner, err := p.projectMemberRepo.FindProjectOwnerByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	ownerInfo, err := p.userRepo.FindByID(ctx, owner.UserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	var profileUrl = ownerInfo.DefaultProfileUrl
	if ownerInfo.UploadedProfileUrl != nil {
		profileUrl = *ownerInfo.UploadedProfileUrl
	}

	return &responses.GetProjectDetailResponse{
		ID:                   project.ID.Hex(),
		WorkspaceID:          project.WorkspaceID.Hex(),
		Name:                 project.Name,
		ProjectPrefix:        project.ProjectPrefix,
		Description:          project.Description,
		Status:               project.Status.String(),
		OwnerUserID:          owner.UserID.Hex(),
		OwnerProjectMemberID: owner.ID.Hex(),
		OwnerDisplayName:     ownerInfo.DisplayName,
		OwnerProfileUrl:      profileUrl,
		Positions:            project.Positions,
		Workflows:            project.Workflows,
		AttributeTemplates:   project.AttributeTemplates,
		CreatedAt:            project.CreatedAt,
		CreatedBy:            project.CreatedBy.Hex(),
		UpdatedAt:            project.UpdatedAt,
		UpdatedBy:            project.UpdatedBy.Hex(),
	}, nil
}

func (p *projectServiceImpl) UpdateSetupStatus(ctx context.Context, req *requests.UpdateSetupStatusRequest, userID string) (*models.Project, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if !models.ProjectStatus(req.Status).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidProjectSetupStatus, errutils.BadRequest).WithDebugMessage("Invalid project setup status")
	}

	// Check if the user is owner of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.BadRequest).WithDebugMessage("User not found")
	} else if member.Role != models.ProjectMemberRoleOwner {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("Requester is not owner of the project")
	}

	project, err := p.projectRepo.UpdateSetupStatus(ctx, &repositories.UpdateProjectSetupStatus{
		ProjectID:   bsonProjectID,
		SetupStatus: models.ProjectSetupStatus(req.Status),
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return project, nil
}

func (p *projectServiceImpl) UpdatePositions(ctx context.Context, req *requests.UpdatePositionsRequest, userID string) (*responses.UpdatePositionsResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if len(req.Titles) == 0 {
		return nil, errutils.NewError(exceptions.ErrNoPositionProvided, errutils.BadRequest).WithDebugMessage("No position provided")
	}

	// Check if the user is owner or moderator of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.BadRequest)
	} else if member.Role != models.ProjectMemberRoleOwner && member.Role != models.ProjectMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	currentPositions, err := p.projectRepo.FindPositionByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	deletedPositions := make([]string, 0)
	for _, position := range currentPositions {
		if !array.ContainAny(req.Titles, []string{position}) {
			deletedPositions = append(deletedPositions, position)
		}
	}

	// Check if the deleted positions are used by any member
	members, err := p.projectMemberRepo.FindByProjectIDAndPositions(ctx, bsonProjectID, deletedPositions)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if len(members) > 0 {
		positionMap := make(map[string]struct{})
		errFields := make([]string, 0)
		for _, member := range members {
			if _, exists := positionMap[member.Position]; !exists {
				positionMap[member.Position] = struct{}{}
				errFields = append(errFields, member.Position)
			}
		}

		return nil, errutils.NewError(exceptions.ErrPositionUsedByMember, errutils.BadRequest).WithDebugMessage("Position is used by member").WithFields(errFields...)
	}

	err = p.projectRepo.UpdatePositions(ctx, bsonProjectID, req.Titles)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.UpdatePositionsResponse{
		Message: "Position updated successfully",
	}, nil
}

func (p *projectServiceImpl) ListPositions(ctx context.Context, req *requests.ListPositionsPathParams) ([]string, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound).WithDebugMessage("Project not found")
	}

	positions, err := p.projectRepo.FindPositionByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return positions, nil
}

func (p *projectServiceImpl) AddMembers(ctx context.Context, req *requests.AddProjectMembersRequest, userID string) (*responses.AddProjectMembersResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if len(req.Members) == 0 {
		return &responses.AddProjectMembersResponse{
			Message: "No member added",
		}, nil
	}

	// Check if the user is owner or moderator of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.BadRequest)
	} else if member.Role != models.ProjectMemberRoleOwner && member.Role != models.ProjectMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound)
	}

	createProjMemberReq := make([]repositories.CreateProjectMemberRequest, 0)
	for _, member := range req.Members {
		bsonMemberID, err := bson.ObjectIDFromHex(member.UserID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		// Check if the member already exists
		existingMember, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonMemberID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
		} else if existingMember != nil {
			continue
		}

		createProjMemberReq = append(createProjMemberReq, repositories.CreateProjectMemberRequest{
			UserID:    bsonMemberID,
			ProjectID: bsonProjectID,
			Role:      models.ProjectMemberRole(member.Role),
			Position:  member.Position,
		})
	}

	if len(createProjMemberReq) == 0 {
		return &responses.AddProjectMembersResponse{
			Message: "No member added",
		}, nil
	}

	err = p.projectMemberRepo.CreateMany(ctx, createProjMemberReq)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.AddProjectMembersResponse{
		Message: "Members added successfully",
	}, nil
}

func validateListMembersPaginationRequestSortBy(sortBy string) bool {
	switch sortBy {
	case constant.ProjectMemberFieldDisplayName, constant.ProjectMemberFieldJoinedAt:
		return true
	}
	return false
}

func normalizeListMembersPaginationRequest(req *requests.ListProjectMembersRequest) {
	if req.PaginationRequest.Page <= 0 {
		req.PaginationRequest.Page = 1
	}
	if req.PaginationRequest.PageSize <= 0 {
		req.PaginationRequest.PageSize = 100
	}
	if req.PaginationRequest.SortBy == "" || !validateListMembersPaginationRequestSortBy(req.PaginationRequest.SortBy) {
		req.PaginationRequest.SortBy = constant.ProjectMemberFieldDisplayName
	}
	if req.PaginationRequest.Order == "" {
		req.PaginationRequest.Order = constant.ASC
	}
}

func (p *projectServiceImpl) ListMembers(ctx context.Context, req *requests.ListProjectMembersRequest) (*responses.ListProjectMembersResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	normalizeListMembersPaginationRequest(req)

	projectMembers, err := p.projectMemberRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	var userIDs []bson.ObjectID
	for _, member := range projectMembers {
		userIDs = append(userIDs, member.UserID)
	}

	users, totalUser, err := p.userRepo.SearchWithUserIDs(ctx, &repositories.SearchUserWithUserIDsRequest{
		UserIDs: userIDs,
		Keyword: req.Keyword,
		PaginationRequest: repositories.PaginationRequest{
			Page:     req.PaginationRequest.Page,
			PageSize: req.PaginationRequest.PageSize,
			Order:    req.PaginationRequest.Order,
			SortBy:   req.PaginationRequest.SortBy,
		},
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	// Map projectMembers and members data to response
	userMap := make(map[bson.ObjectID]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	memberResp := make([]responses.ListProjectMembersResponseMember, 0)
	for _, member := range projectMembers {
		if user, exists := userMap[member.UserID]; exists {
			var profileUrl = user.DefaultProfileUrl
			if user.UploadedProfileUrl != nil {
				profileUrl = *user.UploadedProfileUrl
			}

			memberResp = append(memberResp, responses.ListProjectMembersResponseMember{
				UserID:      user.ID.Hex(),
				Email:       user.Email,
				FullName:    user.FullName,
				DisplayName: user.DisplayName,
				ProfileUrl:  profileUrl,
				Role:        member.Role.String(),
				Position:    member.Position,
				JoinedAt:    member.JoinedAt,
			})
		}
	}

	return &responses.ListProjectMembersResponse{
		Members: memberResp,
		PaginationResponse: &responses.PaginationResponse{
			Page:      req.PaginationRequest.Page,
			PageSize:  req.PaginationRequest.PageSize,
			TotalPage: int(math.Ceil(float64(totalUser) / float64(req.PaginationRequest.PageSize))),
			TotalItem: int(totalUser),
		},
	}, nil
}

func (p *projectServiceImpl) UpdateWorkflows(ctx context.Context, req *requests.UpdateWorkflowsRequest, userID string) (*responses.UpdateWorkflowsResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if len(req.Workflows) == 0 {
		return nil, errutils.NewError(exceptions.ErrNoWorkflowProvided, errutils.BadRequest).WithDebugMessage("No workflow provided")
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound)
	}

	// Check if the user is owner or moderator of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	} else if member.Role != models.ProjectMemberRoleOwner && member.Role != models.ProjectMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	currentWorkflows, err := p.projectRepo.FindWorkflowByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	inputtedStatus := make([]string, 0)
	for _, workflow := range req.Workflows {
		inputtedStatus = append(inputtedStatus, workflow.Status)
	}

	deletedWorkflows := make([]string, 0)
	for _, workflow := range currentWorkflows {
		if !array.ContainAny(inputtedStatus, []string{workflow.Status}) {
			deletedWorkflows = append(deletedWorkflows, workflow.Status)
		}
	}

	// Check if the deleted workflows are used by any task
	tasks, err := p.taskRepo.FindByProjectIDAndStatuses(ctx, bsonProjectID, deletedWorkflows)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if len(tasks) > 0 {
		statusMap := make(map[string]struct{})
		errFields := make([]string, 0)
		for _, task := range tasks {
			if _, exists := statusMap[task.Status]; !exists {
				statusMap[task.Status] = struct{}{}
				errFields = append(errFields, task.Status)
			}
		}

		return nil, errutils.NewError(exceptions.ErrWorkflowUsedByTask, errutils.BadRequest).WithDebugMessage("Workflow is used by task").WithFields(errFields...)
	}

	var (
		workflows          []models.ProjectWorkflow
		isDefaultWorkflows []models.ProjectWorkflow
		isDoneWorkflows    []models.ProjectWorkflow
	)
	for _, workflow := range req.Workflows {
		wf := models.ProjectWorkflow{
			Status:           workflow.Status,
			PreviousStatuses: workflow.PreviousStatuses,
			IsDefault:        workflow.IsDefault,
			IsDone:           workflow.IsDone,
		}
		workflows = append(workflows, wf)

		if workflow.IsDefault {
			isDefaultWorkflows = append(isDefaultWorkflows, wf)
		}

		if workflow.IsDone {
			isDoneWorkflows = append(isDoneWorkflows, wf)
		}
	}

	if len(isDefaultWorkflows) == 0 {
		return nil, errutils.NewError(exceptions.ErrNoDefaultWorkflow, errutils.BadRequest).WithDebugMessage("No default workflow")
	} else if len(isDefaultWorkflows) > 1 {
		errFields := make([]string, 0)
		for _, wf := range isDefaultWorkflows {
			errFields = append(errFields, wf.Status)
		}
		return nil, errutils.NewError(exceptions.ErrMultipleDefaultWorkflow, errutils.BadRequest).WithDebugMessage("Multiple default workflow").WithFields(errFields...)
	}

	if len(isDoneWorkflows) == 0 {
		return nil, errutils.NewError(exceptions.ErrNoIsDoneWorkflow, errutils.BadRequest).WithDebugMessage("No is done workflow")
	}

	err = p.projectRepo.UpdateWorkflows(ctx, bsonProjectID, workflows)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.UpdateWorkflowsResponse{
		Message: "Workflow added successfully",
	}, nil
}

func (p *projectServiceImpl) ListWorkflows(ctx context.Context, req *requests.ListWorkflowsPathParams) ([]models.ProjectWorkflow, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound).WithDebugMessage("Project not found")
	}

	workflows, err := p.projectRepo.FindWorkflowByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return workflows, nil
}

func (p *projectServiceImpl) UpdateAttributeTemplates(ctx context.Context, req *requests.UpdateAttributeTemplatesRequest, userID string) (*responses.UpdateAttributeTemplatesResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	for _, attributeTemplate := range req.AttributeTemplates {
		if !models.KeyValuePairType(strings.ToUpper(attributeTemplate.Type)).IsValid() {
			return nil, errutils.NewError(exceptions.ErrInvalidAttributeType, errutils.BadRequest).WithDebugMessage("Invalid attribute type")
		}
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound)
	}

	// Check if the user is owner or moderator of the project
	member, err := p.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	} else if member.Role != models.ProjectMemberRoleOwner && member.Role != models.ProjectMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	var attributeTemplates []models.ProjectAttributeTemplate
	for _, attributeTemplate := range req.AttributeTemplates {
		attributeTemplates = append(attributeTemplates, models.ProjectAttributeTemplate{
			Name: attributeTemplate.Name,
			Type: models.KeyValuePairType(strings.ToUpper(attributeTemplate.Type)),
		})
	}

	err = p.projectRepo.UpdateAttributeTemplates(ctx, bsonProjectID, attributeTemplates)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.UpdateAttributeTemplatesResponse{
		Message: "Attribute template updated successfully",
	}, nil
}

func (p *projectServiceImpl) ListAttributeTemplates(ctx context.Context, req *requests.ListAttributeTemplatesPathParams) ([]models.ProjectAttributeTemplate, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound).WithDebugMessage("Project not found")
	}

	attributeTemplates, err := p.projectRepo.FindAttributeTemplatesByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return attributeTemplates, nil
}
