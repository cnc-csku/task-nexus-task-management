package router

import "github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"

type Router struct {
	healthCheck rest.HealthCheckHandler
	common      rest.CommonHandler
	member      rest.MemberHandler
}

func NewRouter(
	healthCheck rest.HealthCheckHandler,
	common rest.CommonHandler,
	member rest.MemberHandler,
) *Router {
	return &Router{
		healthCheck: healthCheck,
		common:      common,
		member:      member,
	}
}
