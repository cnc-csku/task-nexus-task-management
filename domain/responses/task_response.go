package responses

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetTaskDetailResponse struct {
	ID                 string                             `json:"id"`
	TaskRef            string                             `json:"taskRef"`
	ProjectID          string                             `json:"projectId"`
	Title              string                             `json:"title"`
	Description        string                             `json:"description"`
	ParentID           *bson.ObjectID                     `json:"parentId"`
	Type               models.TaskType                    `json:"type"`
	Status             string                             `json:"status"`
	Priority           *models.TaskPriority               `json:"priority"`
	Approvals          []models.TaskApproval              `json:"approvals"`
	Assignees          []models.TaskAssignee              `json:"assignees"`
	ChildrenPoint      int                                `json:"childrenPoint"`
	HasChildren        bool                               `json:"hasChildren"`
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
