package console

import "errors"

// Exported Errors
var (
	ErrConsoleURLCannotBeEmpty = errors.New("console base URL cannot be empty")
	ErrUAAURLCannotBeEmpty     = errors.New("UAA URL cannot be empty")
	ErrMissingRefreshToken     = errors.New("missing refresh token")
	ErrNotAuthorized           = errors.New("not authorized")
)
