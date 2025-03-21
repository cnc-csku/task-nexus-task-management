package router

import (
	"github.com/labstack/echo/v4"
)

func (r *Router) RegisterAPIRouter(e *echo.Echo) {
	api := e.Group("/api")

	api.GET("/health", r.healthCheck.HealthCheck)

	common := api.Group("/common/v1")
	{
		common.POST("/generate-put-presigned-url", r.common.GeneratePutPresignedURL, r.authMiddleware.Middleware)
	}

	auth := api.Group("/auth/v1")
	{
		auth.POST("/register", r.user.Register)
		auth.POST("/login", r.user.Login)
		auth.GET("/profile", r.user.GetMyProfile, r.authMiddleware.Middleware)
		auth.GET("/search", r.user.SearchUser, r.authMiddleware.Middleware)
		auth.PUT("/profile", r.user.UpdateProfile, r.authMiddleware.Middleware)
	}

	users := api.Group("/users/v1")
	{
		users.GET("/:userId/profile", r.user.GetUserProfile, r.authMiddleware.Middleware)
	}

	workspaces := api.Group("/workspaces/v1")
	{
		workspaces.GET("/own-workspaces", r.workspace.ListOwnWorkspace, r.authMiddleware.Middleware)
		workspaces.GET("/:workspaceId/members", r.workspace.ListWorkspaceMembers, r.authMiddleware.Middleware)
		workspaces.GET("/:workspaceId/my-projects", r.project.ListMyProjects, r.authMiddleware.Middleware)
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
		projects.GET("/:projectId", r.project.GetProjectDetail, r.authMiddleware.Middleware)

		// Setup
		projects.PUT("/:projectId/setup-status", r.project.UpdateSetupStatus, r.authMiddleware.Middleware)

		// Positions
		projects.PUT("/:projectId/positions", r.project.UpdatePositions, r.authMiddleware.Middleware)
		projects.GET("/:projectId/positions", r.project.ListPositions, r.authMiddleware.Middleware)

		// Members
		projects.POST("/:projectId/members", r.project.AddMembers, r.authMiddleware.Middleware)
		projects.GET("/:projectId/members", r.project.ListMembers, r.authMiddleware.Middleware)

		// Workflow
		projects.PUT("/:projectId/workflows", r.project.UpdateWorkflows, r.authMiddleware.Middleware)
		projects.GET("/:projectId/workflows", r.project.ListWorkflows, r.authMiddleware.Middleware)

		// Sprint
		projects.POST("/:projectId/sprints", r.sprint.Create, r.authMiddleware.Middleware)
		projects.GET("/:projectId/sprints/:sprintId", r.sprint.GetByID, r.authMiddleware.Middleware)
		projects.PUT("/:projectId/sprints/:sprintId", r.sprint.Edit, r.authMiddleware.Middleware)
		projects.GET("/:projectId/sprints", r.sprint.List, r.authMiddleware.Middleware)
		projects.PUT("/:projectId/sprints/:currentSprintId/complete", r.sprint.CompleteSprint, r.authMiddleware.Middleware)
		projects.PUT("/:projectId/sprints/:sprintId/status", r.sprint.UpdateStatus, r.authMiddleware.Middleware)
		projects.DELETE("/:projectId/sprints/:sprintId", r.sprint.Delete, r.authMiddleware.Middleware)

		// Attribute Templates
		projects.PUT("/:projectId/attribute-templates", r.project.UpdateAttributeTemplates, r.authMiddleware.Middleware)
		projects.GET("/:projectId/attribute-templates", r.project.ListAttributeTemplates, r.authMiddleware.Middleware)

		// Project Members
		// Position
		projects.PUT("/:projectId/members/position", r.projectMember.UpdatePosition, r.authMiddleware.Middleware)
	}

	tasks := api.Group("/projects/v1/:projectId/tasks/v1")
	{
		tasks.POST("", r.task.Create, r.authMiddleware.Middleware)
		tasks.GET("/:taskRef", r.task.GetTaskDetail, r.authMiddleware.Middleware)

		tasks.GET("/epic", r.task.ListEpicTasks, r.authMiddleware.Middleware)
		tasks.GET("", r.task.SearchTask, r.authMiddleware.Middleware)

		tasks.PUT("/:taskRef/detail", r.task.UpdateDetail, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/title", r.task.UpdateTitle, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/parent", r.task.UpdateParentID, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/type", r.task.UpdateType, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/status", r.task.UpdateStatus, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/approvals", r.task.UpdateApprovals, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/approve", r.task.ApproveTask, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/assignees", r.task.UpdateAssignees, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/sprint", r.task.UpdateSprint, r.authMiddleware.Middleware)
		tasks.PUT("/:taskRef/attributes", r.task.UpdateAttributes, r.authMiddleware.Middleware)

		tasks.POST("/:taskRef/comments", r.taskComment.Create, r.authMiddleware.Middleware)
		tasks.GET("/:taskRef/comments", r.taskComment.List, r.authMiddleware.Middleware)

		// llm
		tasks.GET("/:taskRef/generate-description", r.task.GenerateDescription, r.authMiddleware.Middleware)
	}

	reports := api.Group("/projects/v1/:projectId/reports/v1")
	{
		reports.GET("/status-overview", r.report.GetStatusOverview, r.authMiddleware.Middleware)
		reports.GET("/priority-overview", r.report.GetPriorityOverview, r.authMiddleware.Middleware)
		reports.GET("/type-overview", r.report.GetTypeOverview, r.authMiddleware.Middleware)
		reports.GET("/assignee-overview", r.report.GetAssigneeOverview, r.authMiddleware.Middleware)
		reports.GET("/epic-task-overview", r.report.GetEpicTaskOverview, r.authMiddleware.Middleware)
	}

	setup := api.Group("/setup/v1")
	{
		setup.GET("", r.common.GetSetupStatus)
		setup.POST("/workspace", r.workspace.SetupWorkspace, r.authMiddleware.Middleware)
		setup.POST("/user", r.user.SetupUser)
	}
}
