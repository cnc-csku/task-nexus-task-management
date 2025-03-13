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

type TaskHandler interface {
	Create(c echo.Context) error
	GetTaskDetail(c echo.Context) error
	ListEpicTasks(c echo.Context) error
	SearchTask(c echo.Context) error
	UpdateDetail(c echo.Context) error
	UpdateStatus(c echo.Context) error
	UpdateApprovals(c echo.Context) error
	ApproveTask(c echo.Context) error
	UpdateAssignees(c echo.Context) error
	UpdateSprint(c echo.Context) error
	UpdateAttributes(c echo.Context) error
}

type taskHandlerImpl struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) TaskHandler {
	return &taskHandlerImpl{
		taskService: taskService,
	}
}

func (h *taskHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) GetTaskDetail(c echo.Context) error {
	req := new(requests.GetTaskDetailPathParam)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.GetTaskDetail(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) ListEpicTasks(c echo.Context) error {
	req := new(requests.ListEpicTasksPathParam)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.ListEpicTasks(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) SearchTask(c echo.Context) error {
	req := new(requests.SearchTaskParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.SearchTask(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateDetail(c echo.Context) error {
	req := new(requests.UpdateTaskDetailRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateDetail(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateStatus(c echo.Context) error {
	req := new(requests.UpdateTaskStatusRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateStatus(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateApprovals(c echo.Context) error {
	req := new(requests.UpdateTaskApprovalsRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateApprovals(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) ApproveTask(c echo.Context) error {
	req := new(requests.ApproveTaskRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.ApproveTask(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateAssignees(c echo.Context) error {
	req := new(requests.UpdateTaskAssigneesRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateAssignees(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateSprint(c echo.Context) error {
	req := new(requests.UpdateTaskSprintRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateSprint(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *taskHandlerImpl) UpdateAttributes(c echo.Context) error {
	req := new(requests.UpdateTaskAttributesRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := h.taskService.UpdateAttributes(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}
