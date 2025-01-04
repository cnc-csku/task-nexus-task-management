package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type CommonHandler interface {
	TestNotification(c echo.Context) error
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

func (h *commonHandler) TestNotification(c echo.Context) error {
	req := new(requests.TestNotificationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &echo.Map{
			"message": "Invalid request",
		})
	}

	res, err := h.commonService.TestNotification(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &echo.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, res)
}
