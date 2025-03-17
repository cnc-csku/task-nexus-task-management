package requests

type CreateTaskCommentRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

type ListTaskCommentPathParams struct {
	ProjectID string `param:"projectId" validate:"required"`
	TaskRef   string `param:"taskRef" validate:"required"`
}
