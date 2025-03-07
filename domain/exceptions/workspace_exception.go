package exceptions

import "github.com/pkg/errors"

var (
	ErrInvalidWorkspaceID           = errors.New("invalid workspace ID")
	ErrMemberNotFoundInWorkspace    = errors.New("member not found in workspace")
	ErrRequesterNotFoundInWorkspace = errors.New("requester not found in workspace")
	ErrMemberAlreadyInWorkspace     = errors.New("member already in workspace")
)
