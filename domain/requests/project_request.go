package requests

type CreateProjectRequest struct {
	Name          string `json:"name" validate:"required"`
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	ProjectPrefix string `json:"projectPrefix" validate:"required"`
	Description   string `json:"description"`
}

type ListMyProjectsPathParams struct {
	WorkspaceID string `param:"workspaceId" validate:"required"`
}

type GetProjectsDetailPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type UpdatePositionsRequest struct {
	ProjectID string   `param:"projectId" validate:"required"`
	Titles    []string `json:"titles" validate:"required"`
}

type ListPositionsPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type AddProjectMembersRequest struct {
	ProjectID string                           `param:"projectId" validate:"required"`
	Members   []AddProjectMembersRequestMember `json:"members" validate:"required,dive"`
}

type AddProjectMembersRequestMember struct {
	UserID   string `json:"userId" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=MEMBER MODERATOR"`
	Position string `json:"position" validate:"required"`
}

type ListProjectMembersRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	Keyword   string `query:"keyword"`
	PaginationRequest
}

type UpdateWorkflowsRequest struct {
	ProjectID string                           `param:"projectId" validate:"required"`
	Workflows []UpdateWorkflowsRequestWorkflow `json:"workflows" validate:"required,dive"`
}

type UpdateWorkflowsRequestWorkflow struct {
	PreviousStatuses []string `json:"previousStatuses"`
	Status           string   `json:"status" validate:"required"`
	IsDefault        bool     `json:"isDefault"`
}

type ListWorkflowsPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type UpdateAttributeTemplatesRequest struct {
	ProjectID          string                                     `param:"projectId" validate:"required"`
	AttributeTemplates []UpdateAttributeTemplatesRequestAttribute `json:"attributesTemplates" validate:"required,dive"`
}

type UpdateAttributeTemplatesRequestAttribute struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

type ListAttributeTemplatesPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
}
