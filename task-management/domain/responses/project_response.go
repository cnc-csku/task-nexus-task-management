package responses

import (
	"time"
)

type CreateProjectResponse struct {
	ID            string  `json:"id"`
	WorkspaceID   string  `json:"workspaceId"`
	Name          string  `json:"name"`
	ProjectPrefix string  `json:"projectPrefix"`
	Description   *string `json:"description"`
}

type AddPositionsResponse struct {
	Message string `json:"message"`
}

type AddProjectMembersResponse struct {
	Message string `json:"message"`
}

type ListProjectMembersResponse struct {
	Members            []ListProjectMembersResponseMember `json:"members"`
	PaginationResponse *PaginationResponse                `json:"paginationResponse"`
}

type ListProjectMembersResponseMember struct {
	UserID      string     `json:"userId"`
	Email       string     `json:"email"`
	FullName    string     `json:"fullName"`
	DisplayName string     `json:"displayName"`
	ProfileUrl  string     `json:"profileUrl"`
	Role        string     `json:"role"`
	Position    string     `json:"position"`
	JoinedAt    time.Time  `json:"joinedAt"`
	RemovedAt   *time.Time `json:"removedAt"`
}

type AddWorkflowsResponse struct {
	Message string `json:"message"`
}
