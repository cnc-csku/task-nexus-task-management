package requests

type CreateProjectRequest struct {
	Name          string   `json:"name" validate:"required"`
	WorkspaceID   string   `json:"workspaceID" validate:"required"`
	ProjectPrefix string   `json:"projectPrefix" validate:"required"`
	Description   string   `json:"description"`
	UserIDs       []string `json:"userIDs"`
}
