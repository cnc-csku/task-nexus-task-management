package requests

type CreateInvitationRequest struct {
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	InviteeEmail  string `json:"inviteeEmail" validate:"required"`
	Role          string `json:"role" validate:"required,oneof=MODERATOR MEMBER"`
	CustomMessage string `json:"customMessage"`
}

type ListInvitationForWorkspaceOwnerParams struct {
	WorkspaceID string `param:"workspaceId" validate:"required"`
	Keyword     string `query:"keyword"`
	SearchBy    string `query:"searchBy"`
	PaginationRequest
}

type UserResponseInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required"`
	Action       string `json:"action" validate:"required"`
}
