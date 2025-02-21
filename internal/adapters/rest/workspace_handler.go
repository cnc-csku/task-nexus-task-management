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

type WorkspaceHandler interface {
	SetupWorkspace(c echo.Context) error
	ListOwnWorkspace(c echo.Context) error
	ListWorkspaceMembers(c echo.Context) error
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

func (w *workspaceHandlerImpl) ListOwnWorkspace(c echo.Context) error {
	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)

	workspaces, err := w.workspaceService.ListOwnWorkspace(c.Request().Context(), userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, workspaces)
}

func (w *workspaceHandlerImpl) ListWorkspaceMembers(c echo.Context) error {
	req := new(requests.ListWorkspaceMemberRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	members, err := w.workspaceService.ListWorkspaceMembers(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, members)
}
