package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type CommonHandler interface {
	GetSetupStatus(c echo.Context) error
}

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

func (ch *commonHandler) GetSetupStatus(c echo.Context) error {
	res, err := ch.commonService.GetSetupStatus(c.Request().Context())
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)

}
