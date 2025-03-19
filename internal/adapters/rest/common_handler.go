package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus-go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type CommonHandler interface {
	GetSetupStatus(c echo.Context) error
	GeneratePutPresignedURL(c echo.Context) error
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

func (ch *commonHandler) GeneratePutPresignedURL(c echo.Context) error {
	req := new(requests.GeneratePutPresignedURLRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := ch.commonService.GeneratePutPresignedURL(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}
