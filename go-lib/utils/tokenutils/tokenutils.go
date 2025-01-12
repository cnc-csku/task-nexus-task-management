package tokenutils

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func GetTokenFromEchoHeader(c echo.Context) (string, error) {
	bearer := c.Request().Header.Get("Authorization")

	bearer = strings.TrimSpace(bearer)
	splittedToken := strings.Split(bearer, "Bearer ")
	if len(splittedToken) != 2 {
		return "", ErrInvalidToken
	}

	token := splittedToken[1]

	return token, nil
}

func GetProfileOnEchoContext(c echo.Context) interface{} {
	profile := c.Get("profile")

	return profile
}
