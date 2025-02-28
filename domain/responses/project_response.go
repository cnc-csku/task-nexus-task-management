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

type ListMyProjectsResponse struct {
	ID                   string    `json:"id"`
	WorkspaceID          string    `json:"workspaceId"`
	Name                 string    `json:"name"`
	ProjectPrefix        string    `json:"projectPrefix"`
	Description          *string   `json:"description"`
	Status               string    `json:"status"`
	OwnerUserID          string    `json:"ownerUserId"`
	OwnerProjectMemberID string    `json:"ownerProjectMemberId"`
	OwnerDisplayName     string    `json:"ownerDisplayName"`
	OwnerProfileUrl      string    `json:"ownerProfileUrl"`
	CreatedAt            time.Time `json:"createdAt"`
	CreatedBy            string    `json:"createdBy"`
	UpdatedAt            time.Time `json:"updatedAt"`
	UpdatedBy            string    `json:"updatedBy"`
}

type GetMyProjectDetailResponse struct {
	ID                   string    `json:"id"`
	WorkspaceID          string    `json:"workspaceId"`
	Name                 string    `json:"name"`
	ProjectPrefix        string    `json:"projectPrefix"`
	Description          *string   `json:"description"`
	Status               string    `json:"status"`
	OwnerUserID          string    `json:"ownerUserId"`
	OwnerProjectMemberID string    `json:"ownerProjectMemberId"`
	OwnerDisplayName     string    `json:"ownerDisplayName"`
	OwnerProfileUrl      string    `json:"ownerProfileUrl"`
	CreatedAt            time.Time `json:"createdAt"`
	CreatedBy            string    `json:"createdBy"`
	UpdatedAt            time.Time `json:"updatedAt"`
	UpdatedBy            string    `json:"updatedBy"`
}

type UpdatePositionsResponse struct {
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

type UpdateWorkflowsResponse struct {
	Message string `json:"message"`
}

type UpdateAttributeTemplatesResponse struct {
	Message string `json:"message"`
}
