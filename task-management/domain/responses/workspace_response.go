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
