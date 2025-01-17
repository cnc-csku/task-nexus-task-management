package exceptions

import "github.com/pkg/errors"

var (
	ErrInvalidWorkspaceID        = errors.New("invalid workspace ID")
	ErrMemberNotFoundInWorkspace = errors.New("member not found in workspace")
	ErrMemberAlreadyInWorkspace  = errors.New("member already in workspace")
)
