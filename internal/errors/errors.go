package app_errors

import "errors"

var (
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInternalServer      = errors.New("internal server error")
)
