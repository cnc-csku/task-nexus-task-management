package exceptions

import "github.com/pkg/errors"

var (
	ErrProjectNameAlreadyExists   = errors.New("project name already exists")
	ErrProjectPrefixAlreadyExists = errors.New("project prefix already exists")
	ErrProjectNotFound            = errors.New("project not found")
	ErrDefaultWorkflowNotFound    = errors.New("default workflow not found")
)
