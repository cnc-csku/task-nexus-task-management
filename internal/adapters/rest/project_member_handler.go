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

type ProjectMemberHandler interface {
	UpdatePosition(c echo.Context) error
}

type projectMemberHandlerImpl struct {
	projectMemberService services.ProjectMemberService
}

func NewProjectMemberHandler(projectMemberService services.ProjectMemberService) ProjectMemberHandler {
	return &projectMemberHandlerImpl{
		projectMemberService: projectMemberService,
	}
}

func (u *projectMemberHandlerImpl) UpdatePosition(c echo.Context) error {
	req := new(requests.UpdateMemberPositionRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	projectMember, err := u.projectMemberService.UpdatePosition(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, projectMember)
}
