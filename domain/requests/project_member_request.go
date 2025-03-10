package requests

type UpdateMemberPositionRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
	Position  string `json:"position" validate:"required"`
}
