package exceptions

import "github.com/pkg/errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInternalError     = errors.New("internal server error")
	ErrInvalidReqPayload = errors.New("invalid request payload")
	ErrPermissionDenied  = errors.New("permission denied")
)
