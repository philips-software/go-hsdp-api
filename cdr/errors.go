package cdr

import (
	"errors"
)

// Errors
var (
	ErrCDRURLCannotBeEmpty = errors.New("base CDR URL cannot be empty")
	ErrEmptyResult         = errors.New("empty result")
	ErrMissingAcceptHeader = errors.New("missing accept header")
)
