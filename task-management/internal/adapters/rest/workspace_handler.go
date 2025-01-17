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

type WorkspaceHandler interface {
	SetupWorkspace(c echo.Context) error
}

type workspaceHandlerImpl struct {
	workspaceService services.WorkspaceService
}

func NewWorkspaceHandler(workspaceService services.WorkspaceService) WorkspaceHandler {
	return &workspaceHandlerImpl{
		workspaceService: workspaceService,
	}
}

func (w *workspaceHandlerImpl) SetupWorkspace(c echo.Context) error {
	req := new(requests.CreateWorkspaceRequest)

	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}
	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)

	workspace, err := w.workspaceService.SetupWorkspace(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusCreated, workspace)
}
