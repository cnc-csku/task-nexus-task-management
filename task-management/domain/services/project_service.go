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
}

type projectServiceImpl struct {
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
	config      *config.Config
}

func NewProjectService(userRepo repositories.UserRepository, projectRepo repositories.ProjectRepository, config *config.Config) ProjectService {
	return &projectServiceImpl{
		userRepo:    userRepo,
		projectRepo: projectRepo,
		config:      config,
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

	var users []models.User
	if req.UserIDs != nil {
		var bsonUserIds []bson.ObjectID
		for _, id := range req.UserIDs {
			bsonUserId, err := bson.ObjectIDFromHex(id)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}
			bsonUserIds = append(bsonUserIds, bsonUserId)
		}
		users, err = p.userRepo.FindByIDs(ctx, bsonUserIds)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
		}
	}

	var members []models.Member
	for _, user := range users {
		members = append(members, models.Member{
			UserID:   user.ID,
			FullName: user.FullName,
		})
	}

	project := &repositories.CreateProjectRequest{
		WorkspaceID:   bsonWorkspaceID,
		Name:          req.Name,
		ProjectPrefix: req.ProjectPrefix,
		Description:   req.Description,
		Status:        models.ProjectStatus_ACTIVE,
		Members:       members,
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
