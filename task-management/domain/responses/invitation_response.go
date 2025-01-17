package responses

import "time"

type CreateInvitationResponse struct {
	Message string `json:"message"`
}

type ListInvitationForUserResponse struct {
	Invitations []InvitationForUserResponse `json:"invitations"`
}

type InvitationForUserResponse struct {
	InvitationID       string     `json:"invitationId"`
	WorkspaceID        string     `json:"workspaceId"`
	WorkspaceName      string     `json:"workspaceName"`
	Status             string     `json:"status"`
	CustomMessage      *string    `json:"customMessage"`
	InvitedAt          string     `json:"invitedAt"`
	InviterDisplayName string     `json:"inviterDisplayName"`
	InviterFullName    string     `json:"inviterFullName"`
	InviterUserID      string     `json:"inviterUserId"`
	ExpiredAt          string     `json:"expiredAt"`
	IsExpired          bool       `json:"isExpired"`
	RespondedAt        *time.Time `json:"respondedAt"`
}

type ListInvitationForWorkspaceOwnerResponse struct {
	Invitations        []InvitationForWorkspaceOwnerResponse `json:"invitations"`
	PaginationResponse PaginationResponse           `json:"paginationResponse"`
}

type InvitationForWorkspaceOwnerResponse struct {
	InvitationID       string     `json:"invitationId"`
	WorkspaceID        string     `json:"workspaceId"`
	WorkspaceName      string     `json:"workspaceName"`
	Status             string     `json:"status"`
	CustomMessage      *string    `json:"customMessage"`
	InvitedAt          string     `json:"invitedAt"`
	InviteeDisplayName string     `json:"inviteeDisplayName"`
	InviteeFullName    string     `json:"inviteeFullName"`
	InviteeUserID      string     `json:"inviteeUserId"`
	InviterDisplayName string     `json:"inviterDisplayName"`
	InviterFullName    string     `json:"inviterFullName"`
	InviterUserID      string     `json:"inviterUserId"`
	ExpiredAt          string     `json:"expiredAt"`
	IsExpired          bool       `json:"isExpired"`
	RespondedAt        *time.Time `json:"respondedAt"`
}

type UserResponseInvitationResponse struct {
	Message string `json:"message"`
}
