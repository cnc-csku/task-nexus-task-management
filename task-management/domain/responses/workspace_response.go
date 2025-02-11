package responses

import (
	"time"
)

type ListOwnWorkspaceResponse struct {
	Workspaces []ListOwnWorkspaceResponseWorkspace `json:"workspaces"`
}

type ListOwnWorkspaceResponseWorkspace struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}

type ListWorkspaceMembersResponse struct {
	Members            []ListWorkspaceMembersResponseWorkspaceMember `json:"members"`
	PaginationResponse PaginationResponse                            `json:"paginationResponse"`
}

type ListWorkspaceMembersResponseWorkspaceMember struct {
	WorkspaceMemberID string    `json:"workspaceMemberId"`
	UserID            string    `json:"userId"`
	Role              string    `json:"role"`
	JoinedAt          time.Time `json:"joinedAt"`
	Email             string    `json:"email"`
	FullName          string    `json:"fullName"`
	DisplayName       string    `json:"displayName"`
	ProfileUrl        string    `json:"profileUrl"`
}
