package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
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
	AddPosition(ctx context.Context, req *requests.AddPositionRequest, userID string) (*responses.AddPositionResponse, *errutils.Error)
}

type projectServiceImpl struct {
	userRepo      repositories.UserRepository
	workspaceRepo repositories.WorkspaceRepository
	projectRepo   repositories.ProjectRepository
	config        *config.Config
}

func NewProjectService(userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository, projectRepo repositories.ProjectRepository, config *config.Config) ProjectService {
	return &projectServiceImpl{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		projectRepo:   projectRepo,
		config:        config,
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
	member, err := p.workspaceRepo.FindWorkspaceMemberByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonUserId)
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

	// var users []models.User
	// if req.UserIDs != nil {
	// 	var bsonUserIds []bson.ObjectID
	// 	for _, id := range req.UserIDs {
	// 		bsonUserId, err := bson.ObjectIDFromHex(id)
	// 		if err != nil {
	// 			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	// 		}
	// 		bsonUserIds = append(bsonUserIds, bsonUserId)
	// 	}
	// 	users, err = p.userRepo.FindByIDs(ctx, bsonUserIds)
	// 	if err != nil {
	// 		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	// 	}
	// }

	// var members []models.ProjectMember
	// for _, user := range users {
	// 	members = append(members, models.ProjectMember{
	// 		UserID:      user.ID,
	// 		DisplayName: user.FullName,
	// 		ProfileUrl:  user.ProfileUrl,
	// 	})
	// }

	owner := &models.ProjectMember{
		UserID:      member.UserID,
		DisplayName: member.DisplayName,
		ProfileUrl:  member.ProfileUrl,
		Role:        models.ProjectMemberRoleOwner,
	}

	project := &repositories.CreateProjectRequest{
		WorkspaceID:   bsonWorkspaceID,
		Name:          req.Name,
		ProjectPrefix: req.ProjectPrefix,
		Description:   req.Description,
		Status:        models.ProjectStatusActive,
		Owner:         owner,
		// Members:       members,
		CreatedBy: bsonUserId,
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

	projects, err := p.projectRepo.FindByWorkspaceIDAndUserID(ctx, bsonWorkspaceID, bsonUserID)
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
	member, err := p.projectRepo.FindMemberByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest)
	}

	project, err := p.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.NotFound)
	}

	return project, nil
}

func (p *projectServiceImpl) AddPosition(ctx context.Context, req *requests.AddPositionRequest, userID string) (*responses.AddPositionResponse, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the user is owner or moderator of the project
	member, err := p.projectRepo.FindMemberByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
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
		return &responses.AddPositionResponse{
			Message: "No new position added",
		}, nil
	}

	err = p.projectRepo.AddPosition(ctx, bsonProjectID, newPositions)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	return &responses.AddPositionResponse{
		Message: "Position added successfully",
	}, nil
}
