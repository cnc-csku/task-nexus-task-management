package requests

type CreateInvitationRequest struct {
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	InviteeUserID string `json:"inviteeUserId" validate:"required"`
	CustomMessage string `json:"customMessage"`
}

type UserResponseInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required"`
	Action       string `json:"action" validate:"required"`
}
