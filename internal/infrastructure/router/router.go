package router

import (
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/task-management/middlewares"
)

type Router struct {
	// Handlers
	healthCheck   rest.HealthCheckHandler
	common        rest.CommonHandler
	user          rest.UserHandler
	project       rest.ProjectHandler
	projectMember rest.ProjectMemberHandler
	invitation    rest.InvitationHandler
	workspace     rest.WorkspaceHandler
	sprint        rest.SprintHandler
	task          rest.TaskHandler
	taskComment   rest.TaskCommentHandler

	// Middlewares
	authMiddleware middlewares.AuthMiddleware
}

func NewRouter(
	authMiddleware middlewares.AuthMiddleware,
	healthCheck rest.HealthCheckHandler,
	common rest.CommonHandler,
	user rest.UserHandler,
	project rest.ProjectHandler,
	projectMember rest.ProjectMemberHandler,
	invitation rest.InvitationHandler,
	workspace rest.WorkspaceHandler,
	sprint rest.SprintHandler,
	task rest.TaskHandler,
	taskComment rest.TaskCommentHandler,
) *Router {
	return &Router{
		authMiddleware: authMiddleware,
		healthCheck:    healthCheck,
		common:         common,
		user:           user,
		project:        project,
		projectMember:  projectMember,
		invitation:     invitation,
		workspace:      workspace,
		sprint:         sprint,
		task:           task,
		taskComment:    taskComment,
	}
}
