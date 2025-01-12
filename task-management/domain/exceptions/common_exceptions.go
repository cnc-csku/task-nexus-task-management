package exceptions

import "errors"

var (
	ErrInternalError     = errors.New("internal server error")
	ErrInvalidReqPayload = errors.New("invalid request payload")
)
