package constant

import "time"

const (
	LogFormatJson = "JSON"
	LogFormatText = "TEXT"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

const (
	InvitationExpirationIn = 7 * (24 * time.Hour) // 7 days
)

const (
	TimeFormat = time.RFC3339
)

// Service constants
const (
	InvitationActionAccept  = "ACCEPT"
	InvitationActionDecline = "DECLINE"
)

const (
	SearchTaskParamsTaskBacklog          = "BACKLOG" // WITH_NO_SPRINT
	SearchTaskParamsTaskWithNoEpicFilter = "WITH_NO_EPIC"
)

// Field names
const (
	UserFieldEmail       = "email"
	UserFieldFullName    = "full_name"
	UserFieldDisplayName = "display_name"
)

const (
	InvitationFieldStatus    = "status"
	InvitationFieldCreatedAt = "created_at"
)

const (
	ProjectMemberFieldDisplayName = "display_name"
	ProjectMemberFieldJoinedAt    = "joined_at"
)

// File Category
const (
	UserProfileFileCategory = "USER_PROFILE"
)

const (
	UserProfileFileCategoryPath = "user-profile"
)
