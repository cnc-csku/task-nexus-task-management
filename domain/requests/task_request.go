package requests

import (
	"time"
)

type CreateTaskRequest struct {
	ProjectID        string                      `param:"projectId" validate:"required"`
	Title            string                      `json:"title" validate:"required"`
	Description      string                      `json:"description"`
	ParentID         *string                     `json:"parentId"`
	Type             string                      `json:"type" validate:"required"`
	Priority         *string                     `json:"priority"`
	SprintID         *string                     `json:"sprintId"`
	StartDate        *time.Time                  `json:"startDate"`
	DueDate          *time.Time                  `json:"dueDate"`
	Assignees        []CreateTaskRequestAssignee `json:"assignees"`
	ApprovalUserIDs  []string                    `json:"approvalUserIds"`
	AdditionalFields map[string]any              `json:"additionalFields"`
}

type CreateTaskRequestAssignee struct {
	UserID   *string `json:"userId"`
	Position string  `json:"position" validate:"required"`
	Point    *int    `json:"point"`
}

type GetTaskDetailPathParam struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
}

type GetManyTaskDetailParams struct {
	ProjectID string   `param:"projectId" validate:"required"`
	TaskRefs  []string `query:"taskRefs" validate:"required"`
}

type ListEpicTasksPathParam struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type SearchTaskParams struct {
	ProjectID       string   `param:"projectId" validate:"required"`
	SprintIDs       []string `query:"sprintIds"`
	IsTaskInBacklog *bool    `query:"isTaskInBacklog"` // Task with no sprint
	EpicTaskID      *string  `query:"epicTaskId"`      // Parent_id or WITH_NO_EPIC
	UserIDs         []string `query:"userIds"`
	Positions       []string `query:"positions"`
	Statuses        []string `query:"statuses"`
	SearchKeyword   *string  `query:"searchKeyword"`
}

type GetChildrenTasksParams struct {
	ProjectID     string `param:"projectId" validate:"required"`
	ParentTaskRef string `query:"parentTaskRef" validate:"required"`
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

type UpdateTaskTypeRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
	Type      string `json:"type" validate:"required"`
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
	Position string  `json:"position" validate:"required"`
	UserId   *string `json:"userId"`
	Point    *int    `json:"point"`
}

type UpdateTaskSprintRequest struct {
	ProjectID       string  `param:"projectId" validate:"required"`
	TaskRef         string  `param:"taskRef" validate:"required"`
	CurrentSprintID *string `json:"currentSprintId"`
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

type GenerateDescriptionRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
	Prompt    string `json:"prompt"`
}
