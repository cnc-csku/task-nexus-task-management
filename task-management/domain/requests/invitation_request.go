package requests

type CreateInvitationRequest struct {
	WorkspaceID   string `json:"workspaceId" validate:"required"`
	InviteeUserID string `json:"inviteeUserId" validate:"required"`
	Role          string `json:"role" validate:"required,oneof=MODERATOR MEMBER"`
	CustomMessage string `json:"customMessage"`
}

type ListInvitationForWorkspaceOwnerQueryParams struct {
	WorkspaceID       string             `param:"workspaceId" validate:"required"`
	Keyword           string             `query:"keyword"`
	SearchBy          string             `query:"searchBy"`
	PaginationRequest *PaginationRequest `query:"paginationRequest"`
}

type UserResponseInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required"`
	Action       string `json:"action" validate:"required"`
}
