package responses

import (
	"time"
)

type ListOwnWorkspaceResponse struct {
	Workspaces []ListOwnWorkspaceResponseWorkspace `json:"workspaces"`
}

type ListOwnWorkspaceResponseWorkspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
