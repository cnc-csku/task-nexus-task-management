package services

import (
	"context"
	"math"

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
	Create(ctx context.Context, req *requests.CreateProjectRequest, userId string) (*responses.CreateProjectResponse, *errutils.Error)
	ListMyProjects(ctx context.Context, req *requests.ListMyProjectsPathParams, userID string) ([]*models.Project, *errutils.Error)
	GetProjectDetail(ctx context.Context, req *requests.GetProjectsDetailPathParams, userID string) (*models.Project, *errutils.Error)
	AddPositions(ctx context.Context, req *requests.AddPositionsRequest, userID string) (*responses.AddPositionsResponse, *errutils.Error)
	ListPositions(ctx context.Context, req *requests.ListPositionsPathParams) ([]string, *errutils.Error)
	AddMembers(ctx context.Context, req *requests.AddProjectMembersRequest, userID string) (*responses.AddProjectMembersResponse, *errutils.Error)
	ListMembers(ctx context.Context, req *requests.ListProjectMembersRequest) (*responses.ListProjectMembersResponse, *errutils.Error)
	AddWorkflows(ctx context.Context, req *requests.AddWorkflowsRequest, userID string) (*responses.AddWorkflowsResponse, *errutils.Error)
	ListWorkflows(ctx context.Context, req *requests.ListWorkflowsPathParams) ([]models.Workflow, *errutils.Error)
	AddAttributeTemplates(ctx context.Context, req *requests.AddAttributeTemplatesRequest, userID string) (*responses.AddAttributeTemplatesResponse, *errutils.Error)
	ListAttributeTemplates(ctx context.Context, req *requests.ListAttributeTemplatesPathParams) ([]models.AttributeTemplate, *errutils.Error)
}

type projectServiceImpl struct {
	userRepo            repositories.UserRepository
	workspaceRepo       repositories.WorkspaceRepository
	workspaceMemberRepo repositories.WorkspaceMemberRepository
	projectRepo         repositories.ProjectRepository
	projectMemberRepo   repositories.ProjectMemberRepository
	config              *config.Config
}

func NewProjectService(
	userRepo repositories.UserRepository,
	workspaceRepo repositories.WorkspaceRepository,
	workspaceMemberRepo repositories.WorkspaceMemberRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	config *config.Config,
) ProjectService {
	return &projectServiceImpl{
		userRepo:            userRepo,
		workspaceRepo:       workspaceRepo,
		workspaceMemberRepo: workspaceMemberRepo,
		projectRepo:         projectRepo,
		projectMemberRepo:   projectMemberRepo,
		config:              config,
	}
}

func (p *projectServiceImpl) Create(ctx context.Context, req *requests.CreateProjectRequest, userId string) (*responses.CreateProjectResponse, *errutils.Error) {
	bsonUserId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonWorkspaceID, err := bson.ObjectIDFromHex(req.WorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInvalidWorkspaceID, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the creator is owner or moderator of the workspace
	member, err := p.workspaceMemberRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonUserId)
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
		WorkspaceID:   bsonWorkspaceID,
		Name:          req.Name,
		ProjectPrefix: req.ProjectPrefix,
		Description:   req.Description,
		Status:        models.ProjectStatusActive,
		Owner:         owner,
		Workflows:     models.GetDefaultWorkflows(),
		CreatedBy:     bsonUserId,
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

func (p *projectServiceImpl) ListMyProjects(ctx context.Context, req *requests.ListMyProjectsPathParams, userID string) ([]*models.Project, *errutils.Error) {
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

	var projectIDs []bson.ObjectID
	for _, projectMember := range projectMembers {
		projectIDs = append(projectIDs, projectMember.ProjectID)
	}

	projects, err := p.projectRepo.FindByProjectIDsAndWorkspaceID(ctx, projectIDs, bsonWorkspaceID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return projects, nil
}

func (p *projectServiceImpl) GetProjectDetail(ctx context.Context, req *requests.GetProjectsDetailPathParams, userID string) (*models.Project, *errutils.Error) {
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

	return project, nil
}

func (p *projectServiceImpl) AddPositions(ctx context.Context, req *requests.AddPositionsRequest, userID string) (*responses.AddPositionsResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
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

	// Check if the position already exists
	existingPositions, err := p.projectRepo.FindPositionByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	positionMap := make(map[string]struct{})
	for _, position := range existingPositions {
		positionMap[position] = struct{}{}
	}

	var newPositions []string
	for _, position := range req.Title {
		if _, ok := positionMap[position]; !ok {
			newPositions = append(newPositions, position)
		}
	}

	if len(newPositions) == 0 {
		return &responses.AddPositionsResponse{
			Message: "No new position added",
		}, nil
	}

	err = p.projectRepo.AddPositions(ctx, bsonProjectID, newPositions)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.AddPositionsResponse{
		Message: "Position added successfully",
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

	err = p.projectMemberRepo.CreateMany(ctx, createProjMemberReq)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return nil, nil
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
			memberResp = append(memberResp, responses.ListProjectMembersResponseMember{
				UserID:      user.ID.Hex(),
				Email:       user.Email,
				FullName:    user.FullName,
				DisplayName: user.DisplayName,
				ProfileUrl:  user.ProfileUrl,
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

func (p *projectServiceImpl) AddWorkflows(ctx context.Context, req *requests.AddWorkflowsRequest, userID string) (*responses.AddWorkflowsResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
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

	// Check if the workflow already exists
	existingWorkflows, err := p.projectRepo.FindWorkflowByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	workflowMap := make(map[string]struct{})
	for _, workflow := range existingWorkflows {
		workflowMap[workflow.Status] = struct{}{}
	}

	var newWorkflows []models.Workflow
	for _, workflow := range req.Workflows {
		if _, ok := workflowMap[workflow.Status]; !ok {
			newWorkflows = append(newWorkflows, models.Workflow{
				Status:           workflow.Status,
				PreviousStatuses: workflow.PreviousStatuses,
			})
		}
	}

	if len(newWorkflows) == 0 {
		return &responses.AddWorkflowsResponse{
			Message: "No new workflow added",
		}, nil
	}

	err = p.projectRepo.AddWorkflows(ctx, bsonProjectID, newWorkflows)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.AddWorkflowsResponse{
		Message: "Workflow added successfully",
	}, nil
}

func (p *projectServiceImpl) ListWorkflows(ctx context.Context, req *requests.ListWorkflowsPathParams) ([]models.Workflow, *errutils.Error) {
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

func (p *projectServiceImpl) AddAttributeTemplates(ctx context.Context, req *requests.AddAttributeTemplatesRequest, userID string) (*responses.AddAttributeTemplatesResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	for _, attributeTemplate := range req.AttributeTemplates {
		if !models.KeyValuePairType(attributeTemplate.Type).IsValid() {
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

	// Check if the attribute template already exists
	existingAttributeTemplates, err := p.projectRepo.FindAttributeTemplatesByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	attributeTemplateMap := make(map[string]struct{})
	for _, attributeTemplate := range existingAttributeTemplates {
		attributeTemplateMap[attributeTemplate.Name] = struct{}{}
	}

	var newAttributeTemplates []models.AttributeTemplate
	for _, attributeTemplate := range req.AttributeTemplates {
		if _, ok := attributeTemplateMap[attributeTemplate.Name]; !ok {
			newAttributeTemplates = append(newAttributeTemplates, models.AttributeTemplate{
				Name: attributeTemplate.Name,
				Type: models.KeyValuePairType(attributeTemplate.Type),
			})
		}
	}

	if len(newAttributeTemplates) == 0 {
		return &responses.AddAttributeTemplatesResponse{
			Message: "No new attribute template added",
		}, nil
	}

	err = p.projectRepo.AddAttributeTemplates(ctx, bsonProjectID, newAttributeTemplates)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.AddAttributeTemplatesResponse{
		Message: "Attribute template added successfully",
	}, nil
}

func (p *projectServiceImpl) ListAttributeTemplates(ctx context.Context, req *requests.ListAttributeTemplatesPathParams) ([]models.AttributeTemplate, *errutils.Error) {
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
