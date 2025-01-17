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

type InvitationHandler interface {
	Create(c echo.Context) error
	ListForUser(c echo.Context) error
	ListForAdmin(c echo.Context) error
	UserResponse(c echo.Context) error
}

type invitationHandlerImpl struct {
	invitationService services.InvitationService
}

func NewInvitationHandler(invitationService services.InvitationService) InvitationHandler {
	return &invitationHandlerImpl{
		invitationService: invitationService,
	}
}

func (u *invitationHandlerImpl) Create(c echo.Context) error {
	req := new(requests.CreateInvitationRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.invitationService.Create(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *invitationHandlerImpl) ListForUser(c echo.Context) error {
	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.invitationService.ListForUser(c.Request().Context(), userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *invitationHandlerImpl) ListForAdmin(c echo.Context) error {
	queryParams := new(requests.ListInvitationForAdminQueryParams)
	if err := c.Bind(queryParams); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.invitationService.ListForAdmin(c.Request().Context(), queryParams, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}

func (u *invitationHandlerImpl) UserResponse(c echo.Context) error {
	req := new(requests.UserResponseInvitationRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	res, err := u.invitationService.UserResponse(c.Request().Context(), req, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, res)
}
