package exceptions

import "github.com/pkg/errors"

var (
	ErrInvitationAlreadySent      = errors.New("invitation already sent")
	ErrInvitationNotFound         = errors.New("invitation not found")
	ErrInvitationAlreadyResponded = errors.New("invitation already responded")
	ErrInvalidInvitationAction    = errors.New("invalid invitation action")
)
