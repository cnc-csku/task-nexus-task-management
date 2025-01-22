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
		auth.GET("/search", r.user.SearchUser, r.authMiddleware.Middleware)
	}

	workspaces := api.Group("/workspaces/v1")
	{
		workspaces.GET("/own-workspaces", r.workspace.ListOwnWorkspace, r.authMiddleware.Middleware)
		workspaces.GET("/:workspaceId/members", r.workspace.ListWorkspaceMembers, r.authMiddleware.Middleware)
	}

	invitations := api.Group("/invitations/v1")
	{
		invitations.POST("", r.invitation.Create, r.authMiddleware.Middleware)
		invitations.GET("/users", r.invitation.ListForUser, r.authMiddleware.Middleware)
		invitations.GET("/:workspaceId/workspaces/owner", r.invitation.ListForWorkspaceOwner, r.authMiddleware.Middleware)
		invitations.PUT("/users", r.invitation.UserResponse, r.authMiddleware.Middleware)
	}

	projects := api.Group("/projects/v1")
	{
		projects.POST("", r.project.Create, r.authMiddleware.Middleware)
		projects.GET("/:workspaceId/my-projects", r.project.ListMyProjects, r.authMiddleware.Middleware)
		projects.GET("/:projectId", r.project.GetProjectDetail, r.authMiddleware.Middleware)

		// Positions
		projects.POST("/:projectId/positions", r.project.AddPositions, r.authMiddleware.Middleware)
		projects.GET("/:projectId/positions", r.project.ListPositions, r.authMiddleware.Middleware)

		// Members
		projects.POST("/:projectId/members", r.project.AddMembers, r.authMiddleware.Middleware)
		projects.GET("/:projectId/members", r.project.ListMembers, r.authMiddleware.Middleware)

		// Workflow
		projects.POST("/:projectId/workflows", r.project.AddWorkflows, r.authMiddleware.Middleware)
		projects.GET("/:projectId/workflows", r.project.ListWorkflows, r.authMiddleware.Middleware)
	}

	setup := api.Group("/setup/v1")
	{
		setup.GET("", r.common.GetSetupStatus)
		setup.POST("/workspace", r.workspace.SetupWorkspace, r.authMiddleware.Middleware)
		setup.POST("/user", r.user.SetupUser)
	}
}
