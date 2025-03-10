package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus-go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type SprintHandler interface {
	Create(c echo.Context) error
	GetByID(c echo.Context) error
	Edit(c echo.Context) error
	ListByProjectID(c echo.Context) error
}

type sprintHandlerImpl struct {
	sprintService services.SprintService
}

func NewSprintHandler(sprintService services.SprintService) SprintHandler {
	return &sprintHandlerImpl{
		sprintService: sprintService,
	}
}

func (h *sprintHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateSprintRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	sprint, err := h.sprintService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, sprint)
}

func (h *sprintHandlerImpl) GetByID(c echo.Context) error {
	req := new(requests.GetSprintByIDRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	sprint, err := h.sprintService.GetByID(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, sprint)
}

func (h *sprintHandlerImpl) Edit(c echo.Context) error {
	req := new(requests.EditSprintRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	sprint, err := h.sprintService.Edit(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, sprint)
}

func (h *sprintHandlerImpl) ListByProjectID(c echo.Context) error {
	req := new(requests.ListSprintByProjectIDPathParam)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	sprints, err := h.sprintService.ListByProjectID(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, sprints)
}
