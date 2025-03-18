package responses

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetTaskDetailResponse struct {
	ID                  string                           `json:"id"`
	TaskRef             string                           `json:"taskRef"`
	ProjectID           string                           `json:"projectId"`
	Title               string                           `json:"title"`
	Description         string                           `json:"description"`
	ParentID            *bson.ObjectID                   `json:"parentId"`
	Type                models.TaskType                  `json:"type"`
	Status              string                           `json:"status"`
	Priority            models.TaskPriority              `json:"priority"`
	Approvals           []GetTaskDetailResponseApprovals `json:"approvals"`
	Assignees           []GetTaskDetailResponseAssignee  `json:"assignees"`
	ChildrenPoint       int                              `json:"childrenPoint"`
	HasChildren         bool                             `json:"hasChildren"`
	Sprint              *models.TaskSprint               `json:"sprint"`
	Attributes          []models.TaskAttribute           `json:"attributes"`
	StartDate           *time.Time                       `json:"startDate"`
	DueDate             *time.Time                       `json:"dueDate"`
	CreatedAt           time.Time                        `json:"createdAt"`
	ReporterUserID      string                           `json:"reporterUserId"`
	ReporterDisplayName string                           `json:"reporterDisplayName"`
	ReporterProfileUrl  string                           `json:"reporterProfileUrl"`
	UpdatedAt           time.Time                        `json:"updatedAt"`
	UpdatedBy           string                           `json:"updatedBy"`
	UpdaterDisplayName  string                           `json:"updaterDisplayName"`
}

type GetTaskDetailResponseApprovals struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	ProfileUrl  string `json:"profileUrl"`
	IsApproved  bool   `json:"isApproved"`
	Reason      string `json:"reason"`
}

type GetTaskDetailResponseAssignee struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	ProfileUrl  string `json:"profileUrl"`
	Position    string `json:"position"`
	Point       *int   `json:"point"`
}

type SearchTaskResponse struct {
	ID            string                       `json:"id"`
	TaskRef       string                       `json:"taskRef"`
	Title         string                       `json:"title"`
	ParentID      *string                      `json:"parentId"`
	ParentTitle   *string                      `json:"parentTitle"`
	Type          string                       `json:"type"`
	Status        string                       `json:"status"`
	Assignees     []SearchTaskResponseAssignee `json:"assignees"`
	Approvals     []models.TaskApproval        `json:"approvals"`
	ChildrenPoint int                          `json:"childrenPoint"`
	HasChildren   bool                         `json:"hasChildren"`
	Sprint        *models.TaskSprint           `json:"sprint"`
}

type SearchTaskResponseAssignee struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	ProfileUrl  string `json:"profileUrl"`
	Position    string `json:"position"`
	Point       *int   `json:"point"`
}

type GenerateDescriptionResponse struct {
	Description []*genai.Content `json:"description"`
}
