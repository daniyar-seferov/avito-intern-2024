package app_errors

import "errors"

var (
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInternalServer      = errors.New("internal server error")
	ErrUserPermissions     = errors.New("user does not have the necessary permissions to complete this action")
	ErrTenderId            = errors.New("invalid tender id")
)
