package requests

type CreateWorkspaceRequest struct {
	Name string `json:"name" validate:"required"`
}
