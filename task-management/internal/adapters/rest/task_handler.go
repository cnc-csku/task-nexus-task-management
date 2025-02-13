package rest

import (
	"net/http"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/labstack/echo/v4"
)

type TaskHandler interface {
	Create(c echo.Context) error
	GetTaskDetail(c echo.Context) error
}

type taskHandlerImpl struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) TaskHandler {
	return &taskHandlerImpl{
		taskService: taskService,
	}
}

func (u *taskHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := u.taskService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}

func (u *taskHandlerImpl) GetTaskDetail(c echo.Context) error {
	req := new(requests.GetTaskDetailPathParam)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	resp, err := u.taskService.GetTaskDetail(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, resp)
}
