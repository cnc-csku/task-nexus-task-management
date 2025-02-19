package requests

type CreateWorkspaceRequest struct {
	Name string `json:"name" validate:"required"`
}

type ListWorkspaceMemberRequest struct {
	WorkspaceID string `param:"workspaceId" validate:"required"`
	Keyword     string `json:"keyword"`
	PaginationRequest
}
