package responses

import "github.com/cnc-csku/task-nexus/task-management/domain/models"

type CreateProjectResponse struct {
	ID            string `json:"id"`
	WorkspaceID   string `json:"workspaceId"`
	Name          string `json:"name"`
	ProjectPrefix string `json:"projectPrefix"`
	Description   string `json:"description"`
}

type AddPositionsResponse struct {
	Message string `json:"message"`
}

type AddProjectMembersResponse struct {
	Message string `json:"message"`
}

type ListProjectMembersResponse struct {
	Members            []models.ProjectMember `json:"members"`
	PaginationResponse *PaginationResponse    `json:"paginationResponse"`
}

type AddWorkflowsResponse struct {
	Message string `json:"message"`
}
