package exceptions

import "github.com/pkg/errors"

var (
	ErrInternalError       = errors.New("internal server error")
	ErrInvalidReqPayload   = errors.New("invalid request payload")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrInvalidFileCategory = errors.New("invalid file category")
)
