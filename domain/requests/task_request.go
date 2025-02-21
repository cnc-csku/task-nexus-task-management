package requests

type CreateTaskRequest struct {
	ProjectID   string  `json:"projectId" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	ParentID    *string `json:"parentId"`
	Type        string  `json:"type" validate:"required"`
	SprintID    *string `json:"sprintId"`
}

type GetTaskDetailPathParam struct {
	TaskID string `param:"taskId" validate:"required"`
}
