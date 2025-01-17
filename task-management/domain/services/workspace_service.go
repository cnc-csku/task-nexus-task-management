package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkspaceService interface {
	SetupWorkspace(ctx context.Context, req *requests.CreateWorkspaceRequest, userID string) (*models.Workspace, *errutils.Error)
}

type workspaceServiceImpl struct {
	workspaceRepo     repositories.WorkspaceRepository
	globalSettingRepo repositories.GlobalSettingRepository
	userRepo          repositories.UserRepository
}

func NewWorkspaceService(
	workspaceRepo repositories.WorkspaceRepository,
	globalSettingRepo repositories.GlobalSettingRepository,
	userRepo repositories.UserRepository,
) WorkspaceService {
	return &workspaceServiceImpl{
		workspaceRepo:     workspaceRepo,
		globalSettingRepo: globalSettingRepo,
		userRepo:          userRepo,
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
		return nil, errutils.NewError(exceptions.ErrWorkspaceAlreadySetup, errutils.BadRequest)
	}

	// Check is setup admin
	isSetupAdmin, err := w.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupAdmin)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	if isSetupAdmin == nil {
		err := w.globalSettingRepo.Set(ctx, &models.GlobalSetting{
			Key:   constant.GlobalSettingKeyIsSetupAdmin,
			Type:  models.GlobalSettingTypeBool,
			Value: false,
		})

		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
	}

	if !isSetupAdmin.Value.(bool) {
		return nil, errutils.NewError(exceptions.ErrAdminNotSetup, errutils.BadRequest)
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
		Name:     req.Name,
		UserID:   userObjID,
		UserName: user.DisplayName,
	})

	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	// Set is setup complete
	err = w.globalSettingRepo.Set(ctx, &models.GlobalSetting{
		Key:   constant.GlobalSettingKeyIsSetupWorkspace,
		Type:  models.GlobalSettingTypeBool,
		Value: true,
	})
	return workspace, nil
}
