package requests

type CreateInvitationRequest struct {
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	InviteeUserID string `json:"inviteeUserId" validate:"required"`
	CustomMessage string `json:"customMessage"`
}

type ListInvitationForAdminQueryParams struct {
	WorkspaceID       string            `param:"workspaceId" validate:"required"`
	Keyword           string            `query:"keyword"`
	SearchBy          string            `query:"searchBy"`
	PaginationRequest *PaginationRequest `query:"paginationRequest"`
}

type UserResponseInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required"`
	Action       string `json:"action" validate:"required"`
}
