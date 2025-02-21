package exceptions

import "github.com/pkg/errors"

var (
	ErrSprintNotFound = errors.New("sprint not found")
)
