package rest

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
)

type CommonHandler interface{}

type commonHandler struct {
	commonService services.CommonService
}

func NewCommonHandler(
	commonService services.CommonService,
) CommonHandler {
	return &commonHandler{
		commonService: commonService,
	}
}
