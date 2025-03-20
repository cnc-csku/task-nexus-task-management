package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus-go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type CommonHandler interface {
	GetSetupStatus(c echo.Context) error
	GeneratePutPresignedURL(c echo.Context) error
}

type commonHandler struct {
	commonService        services.CommonService
	globalSettingService services.GlobalSettingService
}

func NewCommonHandler(
	commonService services.CommonService,
	globalSettingService services.GlobalSettingService,
) CommonHandler {
	return &commonHandler{
		commonService:        commonService,
		globalSettingService: globalSettingService,
	}
}

func (ch *commonHandler) GetSetupStatus(c echo.Context) error {
	response := &responses.SetupStatusResponse{}

	setupWorkspaceSetting, err := ch.globalSettingService.GetGlobalSettingByKey(c.Request().Context(), constant.GlobalSettingKeyIsSetupWorkspace)
	if err != nil {
		return err.ToEchoError()
	}
	response.IsSetupWorkspace = setupWorkspaceSetting.Value.(bool)

	setupOwnerSetting, err := ch.globalSettingService.GetGlobalSettingByKey(c.Request().Context(), constant.GlobalSettingKeyIsSetupOwner)
	if err != nil {
		return err.ToEchoError()
	}

	response.IsSetupOwner = setupOwnerSetting.Value.(bool)

	return c.JSON(http.StatusOK, response)

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
