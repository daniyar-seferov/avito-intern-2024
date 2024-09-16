package app_errors

import "errors"

var (
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInternalServer      = errors.New("internal server error")
	ErrUserPermissions     = errors.New("user does not have the necessary permissions")
	ErrInvalidTenderId     = errors.New("invalid tender id")
	ErrInvalidStatus       = errors.New("invalid status")
)
