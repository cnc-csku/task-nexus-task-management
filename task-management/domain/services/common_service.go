package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
)

type CommonService interface {
	GetSetupStatus(ctx context.Context) (*responses.SetupStatusResponse, *errutils.Error)
}

type commonService struct {
	globalSettingRepo repositories.GlobalSettingRepository
}

func NewCommonService(
	globalSettingRepo repositories.GlobalSettingRepository,
) CommonService {
	return &commonService{
		globalSettingRepo: globalSettingRepo,
	}
}

func (c *commonService) GetSetupStatus(ctx context.Context) (*responses.SetupStatusResponse, *errutils.Error) {
	isSetupWorkspace, err := c.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupWorkspace)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	isSetupAdmin, err := c.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupAdmin)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	return &responses.SetupStatusResponse{
		IsSetupWorkspace: isSetupWorkspace.Value.(bool),
		IsSetupAdmin:     isSetupAdmin.Value.(bool),
	}, nil
}