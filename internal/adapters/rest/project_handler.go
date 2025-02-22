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
	AddPositions(c echo.Context) error
	ListPositions(c echo.Context) error
	AddMembers(c echo.Context) error
	ListMembers(c echo.Context) error
	AddWorkflows(c echo.Context) error
	ListWorkflows(c echo.Context) error
	AddAttributeTemplates(c echo.Context) error
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

func (u *projectHandlerImpl) AddPositions(c echo.Context) error {
	req := new(requests.AddPositionsRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.AddPositions(c.Request().Context(), req, userClaims.ID)
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

func (u *projectHandlerImpl) AddWorkflows(c echo.Context) error {
	req := new(requests.AddWorkflowsRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.AddWorkflows(c.Request().Context(), req, userClaims.ID)
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

func (u *projectHandlerImpl) AddAttributeTemplates(c echo.Context) error {
	req := new(requests.AddAttributeTemplatesRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.projectService.AddAttributeTemplates(c.Request().Context(), req, userClaims.ID)
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
