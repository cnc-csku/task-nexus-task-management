package exceptions

import "errors"

var (
	ErrWorkspaceAlreadySetup = errors.New("workspace already setup")
	ErrOwnerAlreadySetup     = errors.New("owner already setup")
	ErrOwnerNotSetup         = errors.New("owner not setup")
)
