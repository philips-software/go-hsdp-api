package cdl

import (
	"errors"
)

// Errors
var (
	ErrCDLURLCannotBeEmpty = errors.New("base CDL URL cannot be empty")
	ErrEmptyResult         = errors.New("empty result")
	ErrCDLForbidden        = errors.New("HTTP 403 CDL response")
	ErrBadRequest          = errors.New("HTTP 400 Bad request")
	ErrNonHttp20xResponse  = errors.New("non HTTP 20x CDL response")
	ErrConflict            = errors.New("HTTP 409 Conflict. Resource/parameter exists already")
)
