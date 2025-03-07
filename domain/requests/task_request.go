package requests

type CreateTaskRequest struct {
	ProjectID   string  `json:"projectId" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
	Type        string  `json:"type" validate:"required"`
	SprintID    *string `json:"sprintId"`
}

type GetTaskDetailPathParam struct {
	TaskRef string `param:"taskRef" validate:"required"`
}

type UpdateTaskDetailRequest struct {
	TaskRef     string  `param:"taskRef" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
	Type        string  `json:"type" validate:"required"`
	Priority    *string `json:"priority"`
}

type UpdateTaskStatusRequest struct {
	TaskID string `param:"taskRef" validate:"required"`
	Status string `json:"status" validate:"required"` // List project's status
}

type UpdateTaskApprovalsRequest struct {
	TaskRef         string   `param:"taskRef" validate:"required"`
	ApprovalUserIDs []string `json:"approvalUserIds" validate:"required"` // List User in the following project
}

type ApproveTaskRequest struct {
	TaskRef string `param:"taskRef" validate:"required"`
	Reason  string `json:"reason"`
}

type UpdateTaskAssigneesRequest struct {
	TaskRef   string                               `param:"taskRef" validate:"required"`
	Assignees []UpdateTaskAssigneesRequestAssignee `json:"assignees" validate:"required,dive"`
}

type UpdateTaskAssigneesRequestAssignee struct {
	Position string `json:"position" validate:"required"` // List project's position
	UserId   string `json:"userId" validate:"required"`   // List User in the following project
	Point    int    `json:"point"`
}

type UpdateTaskSprintRequest struct {
	TaskRef           string   `param:"taskRef" validate:"required"`
	CurrentSprintID   string   `json:"currentSprintId" validate:"required"`
	PreviousSprintIDs []string `json:"previousSprintIds"`
}
