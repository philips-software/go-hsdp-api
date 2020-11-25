package credentials

import (
	"errors"
)

// Exported Errors
var (
	ErrNotFound                       = errors.New("entity not found")
	ErrMissingManagingOrganization    = errors.New("missing managing organization")
	ErrMissingName                    = errors.New("missing name value")
	ErrMissingDescription             = errors.New("missing description value")
	ErrMalformedInputValue            = errors.New("malformed input value")
	ErrMissingOrganization            = errors.New("missing organization")
	ErrMissingProposition             = errors.New("missing proposition")
	ErrMissingGlobalReference         = errors.New("missing global reference")
	ErrMissingProductKey              = errors.New("missing product key")
	ErrNotImplementedByHSDP           = errors.New("method not implemented by HSDP")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrBaseURLCannotBeEmpty           = errors.New("Credentials base URL cannot be empty")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrEmptyResult                    = errors.New("empty result")
)
