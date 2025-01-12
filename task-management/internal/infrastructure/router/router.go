package router

import (
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/task-management/middlewares"
)

type Router struct {
	// Handlers
	healthCheck rest.HealthCheckHandler
	common      rest.CommonHandler
	user        rest.UserHandler

	// Middlewares
	authMiddleware middlewares.AuthMiddleware
}

func NewRouter(
	authMiddleware middlewares.AuthMiddleware,
	healthCheck rest.HealthCheckHandler,
	common rest.CommonHandler,
	user rest.UserHandler,
) *Router {
	return &Router{
		authMiddleware: authMiddleware,
		healthCheck: healthCheck,
		common:      common,
		user:        user,
	}
}
