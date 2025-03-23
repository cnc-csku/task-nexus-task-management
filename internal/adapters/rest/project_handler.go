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

type ProjectHandler interface {
	Create(c echo.Context) error
	ListMyProjects(c echo.Context) error
	GetProjectDetail(c echo.Context) error
	UpdateSetupStatus(c echo.Context) error
	UpdateDetail(c echo.Context) error
	UpdatePositions(c echo.Context) error
	ListPositions(c echo.Context) error
	AddMembers(c echo.Context) error
	ListMembers(c echo.Context) error
	UpdateWorkflows(c echo.Context) error
	ListWorkflows(c echo.Context) error
	UpdateAttributeTemplates(c echo.Context) error
	ListAttributeTemplates(c echo.Context) error
}

type projectHandlerImpl struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) ProjectHandler {
	return &projectHandlerImpl{
		projectService: projectService,
	}
}

func (u *projectHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateProjectRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	project, err := u.projectService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, project)
}

func (u *projectHandlerImpl) ListMyProjects(c echo.Context) error {
	req := new(requests.ListMyProjectsPathParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	projects, err := u.projectService.ListMyProjects(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, projects)
}

func (u *projectHandlerImpl) GetProjectDetail(c echo.Context) error {
	req := new(requests.GetProjectsDetailPathParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	project, err := u.projectService.GetProjectDetail(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, project)
}

func (u *projectHandlerImpl) UpdateSetupStatus(c echo.Context) error {
	req := new(requests.UpdateSetupStatusRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.UpdateSetupStatus(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) UpdateDetail(c echo.Context) error {
	req := new(requests.UpdateProjectDetailRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.UpdateDetail(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) UpdatePositions(c echo.Context) error {
	req := new(requests.UpdatePositionsRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.UpdatePositions(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) ListPositions(c echo.Context) error {
	req := new(requests.ListPositionsPathParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	positions, err := u.projectService.ListPositions(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, positions)
}

func (u *projectHandlerImpl) AddMembers(c echo.Context) error {
	req := new(requests.AddProjectMembersRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.AddMembers(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) ListMembers(c echo.Context) error {
	req := new(requests.ListProjectMembersRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	members, err := u.projectService.ListMembers(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, members)
}

func (u *projectHandlerImpl) UpdateWorkflows(c echo.Context) error {
	req := new(requests.UpdateWorkflowsRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.UpdateWorkflows(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) ListWorkflows(c echo.Context) error {
	req := new(requests.ListWorkflowsPathParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	workflows, err := u.projectService.ListWorkflows(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, workflows)
}

func (u *projectHandlerImpl) UpdateAttributeTemplates(c echo.Context) error {
	req := new(requests.UpdateAttributeTemplatesRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.UpdateAttributeTemplates(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *projectHandlerImpl) ListAttributeTemplates(c echo.Context) error {
	req := new(requests.ListAttributeTemplatesPathParams)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	attributeTemplates, err := u.projectService.ListAttributeTemplates(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, attributeTemplates)
}
