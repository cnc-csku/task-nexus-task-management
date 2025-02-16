package requests

type CreateTaskCommentRequest struct {
	TaskID  string `param:"taskId" validate:"required"`
	Content string `json:"content" validate:"required"`
}
