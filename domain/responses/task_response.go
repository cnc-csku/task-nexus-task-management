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
	Priority           models.TaskPriority                `json:"priority"`
	Approvals          []models.TaskApproval              `json:"approvals"`
	Assignees          []models.TaskAssignee              `json:"assignees"`
	ChildrenPoint      int                                `json:"childrenPoint"`
	HasChildren        bool                               `json:"hasChildren"`
	Sprint             *models.TaskSprint                 `json:"sprint"`
	Attributes         []models.TaskAttribute             `json:"attributes"`
	StartDate          *time.Time                         `json:"startDate"`
	DueDate            *time.Time                         `json:"dueDate"`
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

type SearchTaskResponse struct {
	ID            string                `json:"id"`
	TaskRef       string                `json:"taskRef"`
	Title         string                `json:"title"`
	ParentID      *string               `json:"parentId"`
	ParentTitle   *string               `json:"parentTitle"`
	Type          string                `json:"type"`
	Status        string                `json:"status"`
	Assignees     []models.TaskAssignee `json:"assignees"`
	ChildrenPoint int                   `json:"childrenPoint"`
	HasChildren   bool                  `json:"hasChildren"`
	Sprint        *models.TaskSprint    `json:"sprint"`
}
