package requests

type CreateProjectRequest struct {
	Name          string `json:"name" validate:"required"`
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	ProjectPrefix string `json:"projectPrefix" validate:"required"`
	Description   string `json:"description"`
	// UserIDs       []string `json:"userIds"`
}

type ListMyProjectsPathParams struct {
	WorkspaceID string `param:"workspaceId" validate:"required"`
}

type GetProjectsDetailPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type AddPositionRequest struct {
	ProjectID string   `param:"projectId" validate:"required"`
	Title     []string `json:"title" validate:"required"`
}
