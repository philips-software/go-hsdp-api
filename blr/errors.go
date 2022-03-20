package blr

import (
	"errors"
)

var (
	ErrNotFound                       = errors.New("entity not found")
	ErrBaseURLCannotBeEmpty           = errors.New("base URL cannot be empty")
	ErrEmptyResult                    = errors.New("empty result")
	ErrInvalidEndpointURL             = errors.New("invalid endpoint URL")
	ErrMissingName                    = errors.New("missing name value")
	ErrMissingDescription             = errors.New("missing description value")
	ErrMalformedInputValue            = errors.New("malformed input value")
	ErrMissingOrganization            = errors.New("missing organization")
	ErrMissingProposition             = errors.New("missing proposition")
	ErrMissingGlobalReference         = errors.New("missing global reference")
	ErrNotImplementedByHSDP           = errors.New("method not implemented by HSDP")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
)
