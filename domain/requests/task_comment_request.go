package requests

type CreateTaskCommentRequest struct {
	TaskRef string `param:"taskRef" validate:"required"`
	Content string `json:"content" validate:"required"`
}
