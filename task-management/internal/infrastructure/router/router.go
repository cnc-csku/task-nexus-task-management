package router

import "github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"

type Router struct {
	healthCheck rest.HealthCheckHandler
}

func NewRouter(
	healthCheck rest.HealthCheckHandler,
) *Router {
	return &Router{
		healthCheck: healthCheck,
	}
}
