package audit

import (
	"errors"
)

// Errors
var (
	ErrBaseURLCannotBeEmpty = errors.New("base URL cannot be empty")
	ErrEmptyResult          = errors.New("empty result")
	ErrBadRequest           = errors.New("bad request")
)
