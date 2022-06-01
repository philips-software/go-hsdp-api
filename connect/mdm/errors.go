package mdm

import (
	"errors"
)

var (
	ErrNotFound                       = errors.New("entity not found")
	ErrBaseURLCannotBeEmpty           = errors.New("base URL cannot be empty")
	ErrEmptyResult                    = errors.New("empty result")
	ErrOperationFailed                = errors.New("operation failed")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
)
