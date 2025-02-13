package exceptions

import "github.com/pkg/errors"

var (
	ErrInvalidTaskType       = errors.New("invalid task type")
	ErrParentTaskNotFound    = errors.New("parent task not found")
	ErrInvalidParentTaskType = errors.New("invalid parent task type")
)
