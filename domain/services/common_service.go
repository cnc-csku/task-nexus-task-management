package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/google/uuid"
)

type CommonService interface {
	GetSetupStatus(ctx context.Context) (*responses.SetupStatusResponse, *errutils.Error)
	GeneratePutPresignedURL(ctx context.Context, req *requests.GeneratePutPresignedURLRequest, userID string) (*responses.GeneratePutPresignedURLResponse, *errutils.Error)
}

type commonService struct {
	globalSettingRepo repositories.GlobalSettingRepository
	minioRepo         repositories.MinioRepository
}

func NewCommonService(
	globalSettingRepo repositories.GlobalSettingRepository,
	minioRepo repositories.MinioRepository,
) CommonService {
	return &commonService{
		globalSettingRepo: globalSettingRepo,
		minioRepo:         minioRepo,
	}
}

func (c *commonService) GetSetupStatus(ctx context.Context) (*responses.SetupStatusResponse, *errutils.Error) {
	isSetupWorkspace, err := c.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupWorkspace)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	isSetupOwner, err := c.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupOwner)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	return &responses.SetupStatusResponse{
		IsSetupWorkspace: isSetupWorkspace.Value.(bool),
		IsSetupOwner:     isSetupOwner.Value.(bool),
	}, nil
}

func (c *commonService) GeneratePutPresignedURL(ctx context.Context, req *requests.GeneratePutPresignedURLRequest, userID string) (*responses.GeneratePutPresignedURLResponse, *errutils.Error) {
	fileCategoryPath, err := getFileCategoryPath(req.FileCategory)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInvalidFileCategory, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	key := fmt.Sprintf("%s/%s/%s/%s", fileCategoryPath, userID, uuid.String(), req.FileName)
	url, err := c.minioRepo.GeneratePutPresignedURL(ctx, key)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return &responses.GeneratePutPresignedURLResponse{
		URL: url,
	}, nil
}

func getFileCategoryPath(fileCategory string) (string, error) {
	allowedFileCategories := map[string]string{
		constant.UserProfileFileCategory: constant.UserProfileFileCategoryPath,
	}

	path, exists := allowedFileCategories[fileCategory]
	if !exists {
		return "", errors.New("file category not allowed")
	}

	return path, nil
}
