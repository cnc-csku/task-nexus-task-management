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

type UserHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	GetUserProfile(c echo.Context) error
	SearchUser(c echo.Context) error
	SetupUser(c echo.Context) error
}

type userHandlerImpl struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandlerImpl{
		userService: userService,
	}
}

func (u *userHandlerImpl) Register(c echo.Context) error {
	req := new(requests.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	user, err := u.userService.Register(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, user)
}

func (u *userHandlerImpl) Login(c echo.Context) error {
	req := new(requests.LoginRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	user, err := u.userService.Login(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, user)
}

func (u *userHandlerImpl) GetUserProfile(c echo.Context) error {
	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	user, err := u.userService.FindUserByEmail(c.Request().Context(), userClaims.Email)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, user)
}

func (u *userHandlerImpl) SearchUser(c echo.Context) error {
	params := new(requests.SearchUserParams)
	if err := c.Bind(params); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(params); err != nil {
		return err
	}

	userClaims := tokenutils.GetProfileOnEchoContext(c).(*models.UserCustomClaims)
	users, err := u.userService.Search(c.Request().Context(), params, userClaims.ID)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, users)
}

func (u *userHandlerImpl) SetupUser(c echo.Context) error {
	req := new(requests.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return errutils.NewError(err, errutils.BadRequest).ToEchoError()
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	user, err := u.userService.SetupFirstUser(c.Request().Context(), req)
	if err != nil {
		return err.ToEchoError()
	}

	return c.JSON(http.StatusOK, user)
}
