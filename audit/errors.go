package audit

import (
	"errors"
)

// Errors
var (
	ErrBaseURLCannotBeEmpty           = errors.New("base CDR URL cannot be empty")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrNotImplementedYet              = errors.New("not implemented yet")
	ErrBadRequest                     = errors.New("bad request")
	ErrNonHttp20xResponse             = errors.New("non http 20x Audit response")
)
