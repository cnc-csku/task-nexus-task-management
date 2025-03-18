package exceptions

import "github.com/pkg/errors"

var (
	ErrSprintNotFound        = errors.New("sprint not found")
	ErrInvalidSprintStatus   = errors.New("invalid sprint status")
	ErrDeletedSprintHasTasks = errors.New("deleted sprint has tasks")
)
