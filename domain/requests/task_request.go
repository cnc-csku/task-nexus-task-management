package requests

import (
	"time"
)

type CreateTaskRequest struct {
	ProjectID   string  `param:"projectId" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
	Type        string  `json:"type" validate:"required"`
	SprintID    *string `json:"sprintId"`
}

type GetTaskDetailPathParam struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
}

type ListEpicTasksPathParam struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type SearchTaskParams struct {
	ProjectID     string   `param:"projectId" validate:"required"`
	SprintID      *string  `query:"sprintId"`   // Sprint_id or BACKLOG (WITH_NO_SPRINT)
	EpicTaskID    *string  `query:"epicTaskId"` // Parent_id or WITH_NO_EPIC
	UserIDs       []string `query:"userIds"`
	Positions     []string `query:"positions"`
	Statuses      []string `query:"statuses"`
	SearchKeyword *string  `query:"searchKeyword"`
}

type UpdateTaskDetailRequest struct {
	ProjectID   string     `param:"projectId" validate:"required"`
	TaskRef     string     `param:"taskRef" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" validate:"required"`
	StartDate   *time.Time `json:"startDate"`
	DueDate     *time.Time `json:"dueDate"`
}

type UpdateTaskTitleRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
	Title     string `json:"title" validate:"required"`
}

type UpdateTaskParentIdRequest struct {
	ProjectID string  `param:"projectId" validate:"required"`
	TaskRef   string  `param:"taskRef" validate:"required"`
	ParentID  *string `json:"parentId"`
}

type UpdateTaskStatusRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskID    string `param:"taskRef" validate:"required"`
	Status    string `json:"status" validate:"required"` // List project's status
}

type UpdateTaskApprovalsRequest struct {
	ProjectID       string   `param:"projectId" validate:"required"`
	TaskRef         string   `param:"taskRef" validate:"required"`
	ApprovalUserIDs []string `json:"approvalUserIds" validate:"required"` // List User in the following project
}

type ApproveTaskRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
	Reason    string `json:"reason"`
}

type UpdateTaskAssigneesRequest struct {
	ProjectID string                               `param:"projectId" validate:"required"`
	TaskRef   string                               `param:"taskRef" validate:"required"`
	Assignees []UpdateTaskAssigneesRequestAssignee `json:"assignees" validate:"required,dive"`
}

type UpdateTaskAssigneesRequestAssignee struct {
	ProjectID string `json:"projectId" validate:"required"`
	Position  string `json:"position" validate:"required"` // List project's position
	UserId    string `json:"userId" validate:"required"`   // List User in the following project
	Point     *int   `json:"point"`
}

type UpdateTaskSprintRequest struct {
	ProjectID         string   `param:"projectId" validate:"required"`
	TaskRef           string   `param:"taskRef" validate:"required"`
	CurrentSprintID   string   `json:"currentSprintId" validate:"required"`
	PreviousSprintIDs []string `json:"previousSprintIds"`
}

type UpdateTaskAttributesRequest struct {
	ProjectID  string                                 `param:"projectId" validate:"required"`
	TaskRef    string                                 `param:"taskRef" validate:"required"`
	Attributes []UpdateTaskAttributesRequestAttribute `json:"attributes" validate:"required,dive"`
}

type UpdateTaskAttributesRequestAttribute struct {
	ProjectID string `json:"projectId" validate:"required"`
	Key       string `json:"key" validate:"required"`
	Value     string `json:"value"`
}
