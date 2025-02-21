package middlewares

import (
	"fmt"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus-go-lib/utils/tokenutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	configs *config.Config
}

type AuthMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAdminJWTMiddleware(configs *config.Config) AuthMiddleware {
	return &authMiddleware{
		configs: configs,
	}
}

func (a *authMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := tokenutils.GetTokenFromEchoHeader(c)
		if err != nil {
			return errutils.NewError(err, errutils.Unauthorized).ToEchoError()
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &models.UserCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.configs.JWT.AccessTokenSecret), nil
		})
		if err != nil {
			return errutils.NewError(err, errutils.Unauthorized).ToEchoError()
		}

		// Validate claims
		claims, ok := token.Claims.(*models.UserCustomClaims)
		if !ok || !token.Valid {
			return errutils.NewError(exceptions.ErrInvalidToken, errutils.Unauthorized).ToEchoError()
		}

		// Set claims to context
		c.Set("profile", claims)

		return next(c)
	}
}
