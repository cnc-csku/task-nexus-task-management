package exceptions

import "github.com/pkg/errors"

var (
	ErrProjectNameAlreadyExists   = errors.New("project name already exists")
	ErrProjectPrefixAlreadyExists = errors.New("project prefix already exists")
	ErrProjectNotFound            = errors.New("project not found")
	ErrDefaultWorkflowNotFound    = errors.New("default workflow not found")
	ErrInvalidAttributeType       = errors.New("invalid attribute type")
	ErrNoDefaultWorkflow          = errors.New("no default workflow")
	ErrMultipleDefaultWorkflow    = errors.New("multiple default workflow")
	ErrNoPositionProvided         = errors.New("no position provided")
	ErrNoWorkflowProvided         = errors.New("no workflow provided")
	ErrPositionUsedByMember       = errors.New("position is used by member")
	ErrWorkflowUsedByTask         = errors.New("workflow is used by task")
)
