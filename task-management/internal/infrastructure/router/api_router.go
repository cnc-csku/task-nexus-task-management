package router

import (
	"github.com/labstack/echo/v4"
)

func (r *Router) RegisterAPIRouter(e *echo.Echo) {
	api := e.Group("/api")

	api.GET("/health", r.healthCheck.HealthCheck)

	auth := api.Group("/auth/v1")
	{
		auth.POST("/register", r.user.Register)
		auth.POST("/login", r.user.Login)
		auth.GET("/profile", r.user.GetUserProfile, r.authMiddleware.Middleware)
	}
}
