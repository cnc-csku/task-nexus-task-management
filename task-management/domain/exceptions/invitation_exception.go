package exceptions

import "github.com/pkg/errors"

var (
	ErrInvitationAlreadySent   = errors.New("invitation already sent")
	ErrInvitationNotFound      = errors.New("invitation not found")
	ErrInvalidInvitationStatus = errors.New("invalid invitation status")
	ErrInvalidInvitationAction = errors.New("invalid invitation action")
)
