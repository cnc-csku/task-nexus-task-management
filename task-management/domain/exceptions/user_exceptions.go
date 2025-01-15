package exceptions

import "github.com/pkg/errors"

var (
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrInvalidCredentials         = errors.New("invalid username or password")
	ErrInvalidToken               = errors.New("invalid token")
	ErrUserNotFound               = errors.New("user not found")
	ErrProjectNameAlreadyExists   = errors.New("project name already exists")
	ErrProjectPrefixAlreadyExists = errors.New("project prefix already exists")
	ErrInvalidWorkspaceID         = errors.New("invalid workspace ID")
)
