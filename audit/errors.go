package audit

import (
	"errors"
)

// Errors
var (
	ErrBaseURLCannotBeEmpty = errors.New("base URL cannot be empty")
	ErrEmptyResult          = errors.New("empty result")
	ErrBadRequest           = errors.New("bad request")
	ErrNonHttp20xResponse   = errors.New("non http 20x Audit response")
)
