package exceptions

import "github.com/pkg/errors"

var (
	ErrTaskNotFound                     = errors.New("task not found")
	ErrInvalidTaskType                  = errors.New("invalid task type")
	ErrParentTaskNotFound               = errors.New("parent task not found")
	ErrInvalidParentTaskType            = errors.New("invalid parent task type")
	ErrInvalidTaskPriority              = errors.New("invalid task priority")
	ErrInvalidTaskStatus                = errors.New("invalid task status")
	ErrInvalidSearchTasksSearchFilterBy = errors.New("invalid search tasks search filter by")
	ErrInvalidAttributeKey              = errors.New("invalid attribute key")
)
