package responses

type CreateProjectResponse struct {
	ID            string   `json:"id"`
	WorkspaceID   string   `json:"workspaceId"`
	Name          string   `json:"name"`
	ProjectPrefix string   `json:"projectPrefix"`
	Description   string   `json:"description"`
}
