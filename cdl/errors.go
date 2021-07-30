package cdl

import (
	"errors"
)

// Errors
var (
	ErrCDLURLCannotBeEmpty = errors.New("base CDL URL cannot be empty")
	ErrEmptyResult         = errors.New("empty result")
)
