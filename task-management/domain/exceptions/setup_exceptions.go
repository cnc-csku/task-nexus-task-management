package exceptions

import "errors"

var (
	ErrWorkspaceAlreadySetup = errors.New("workspace already setup")
	ErrAdminAlreadySetup     = errors.New("admin already setup")
	ErrAdminNotSetup         = errors.New("admin not setup")
)
