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

type TaskCommentHandler interface {
	Create(c echo.Context) error
}

type taskCommentHandlerImpl struct {
	taskCommentService services.TaskCommentService
}

func NewTaskCommentHandler(
	taskCommentService services.TaskCommentService,
) TaskCommentHandler {
	return &taskCommentHandlerImpl{
		taskCommentService: taskCommentService,
	}
}

func (h *taskCommentHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateTaskCommentRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	taskComment, err := h.taskCommentService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, taskComment)
}
