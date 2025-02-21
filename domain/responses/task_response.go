package responses

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
)

type GetTaskDetailResponse struct {
	ID                 string                             `json:"id"`
	TaskID             string                             `json:"taskId"`
	ProjectID          string                             `json:"projectId"`
	Title              string                             `json:"title"`
	Description        *string                            `json:"description"`
	ParentID           *string                            `json:"parentId"`
	Type               models.TaskType                    `json:"type"`
	Status             string                             `json:"status"`
	Priority           *models.TaskPriority               `json:"priority"`
	Approval           []models.TaskApproval              `json:"approval"`
	Assignee           []models.TaskAssignee              `json:"assignee"`
	Sprint             *models.TaskSprint                 `json:"sprint"`
	CreatedAt          time.Time                          `json:"createdAt"`
	CreatedBy          string                             `json:"createdBy"`
	CreatorDisplayName string                             `json:"creatorDisplayName"`
	UpdatedAt          time.Time                          `json:"updatedAt"`
	UpdatedBy          string                             `json:"updatedBy"`
	UpdaterDisplayName string                             `json:"updaterDisplayName"`
	TaskComments       []GetTaskDetailResponseTaskComment `json:"taskComments"`
}

type GetTaskDetailResponseTaskComment struct {
	ID              string    `json:"id"`
	Content         string    `json:"content"`
	UserID          string    `json:"userId"`
	UserDisplayName string    `json:"userDisplayName"`
	TaskID          string    `json:"taskId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
