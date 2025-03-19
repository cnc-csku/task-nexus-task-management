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

type SprintHandler interface {
	Create(c echo.Context) error
	GetByID(c echo.Context) error
	Edit(c echo.Context) error
	List(c echo.Context) error
	CompleteSprint(c echo.Context) error
	UpdateStatus(c echo.Context) error
	Delete(c echo.Context) error
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

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	sprint, err := h.sprintService.GetByID(c.Request().Context(), req, userClaims.ID)
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

func (h *sprintHandlerImpl) List(c echo.Context) error {
	req := new(requests.ListSprintPathParam)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	sprints, err := h.sprintService.List(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, sprints)
}

func (h *sprintHandlerImpl) CompleteSprint(c echo.Context) error {
	req := new(requests.CompleteSprintRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.sprintService.CompleteSprint(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *sprintHandlerImpl) UpdateStatus(c echo.Context) error {
	req := new(requests.UpdateSprintStatusRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.sprintService.UpdateStatus(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *sprintHandlerImpl) Delete(c echo.Context) error {
	req := new(requests.DeleteSprintRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.sprintService.Delete(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}
