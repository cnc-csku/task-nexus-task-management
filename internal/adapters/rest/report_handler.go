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
	GetPriorityOverview(c echo.Context) error
	GetTypeOverview(c echo.Context) error
	GetAssigneeOverview(c echo.Context) error
	GetEpicTaskOverview(c echo.Context) error
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
	req := new(requests.GetTaskStatusOverviewRequest)
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

func (h *reportHandlerImpl) GetPriorityOverview(c echo.Context) error {
	req := new(requests.GetTaskPriorityOverviewRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	priorityOverview, err := h.reportService.GetPriorityOverview(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, priorityOverview)
}

func (h *reportHandlerImpl) GetTypeOverview(c echo.Context) error {
	req := new(requests.GetTaskTypeOverviewRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	typeOverview, err := h.reportService.GetTypeOverview(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, typeOverview)
}

func (h *reportHandlerImpl) GetAssigneeOverview(c echo.Context) error {
	req := new(requests.GetTaskAssigneeOverviewRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	assigneeOverview, err := h.reportService.GetAssigneeOverview(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, assigneeOverview)
}

func (h *reportHandlerImpl) GetEpicTaskOverview(c echo.Context) error {
	req := new(requests.GetEpicTaskOverviewRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	epicTaskOverview, err := h.reportService.GetEpicTaskOverview(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, epicTaskOverview)
}
