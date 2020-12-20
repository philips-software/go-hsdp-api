package cdr

import (
	"errors"
)

// Errors
var (
	ErrCDRURLCannotBeEmpty            = errors.New("base CDR URL cannot be empty")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrNotImplementedYet              = errors.New("not implemented yet")
	ErrNonHttp20xResponse             = errors.New("non http 20x console response")
)
