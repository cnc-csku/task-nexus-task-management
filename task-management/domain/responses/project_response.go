package responses

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
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
	UserID      string                   `bson:"user_id" json:"userId"`
	DisplayName string                   `bson:"display_name" json:"displayName"`
	ProfileUrl  string                   `bson:"profile_url" json:"profileUrl"`
	Role        models.ProjectMemberRole `bson:"role" json:"role"`
	Position    string                   `bson:"position" json:"position"`
	JoinedAt    time.Time                `bson:"joined_at" json:"joinedAt"`
	RemovedAt   *time.Time               `bson:"removed_at" json:"removedAt"`
}

type AddWorkflowsResponse struct {
	Message string `json:"message"`
}
