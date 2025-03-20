package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus-go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type ReportHandler interface {
	GetStatusOverview(c echo.Context) error
}

type reportHandlerImpl struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) ReportHandler {
	return &reportHandlerImpl{
		reportService: reportService,
	}
}

func (h *reportHandlerImpl) GetStatusOverview(c echo.Context) error {
	req := new(requests.GetStatusOverviewRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	statusOverview, err := h.reportService.GetStatusOverview(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, statusOverview)
}
