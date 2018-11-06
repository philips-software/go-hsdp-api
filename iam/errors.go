package iam

import "fmt"

// Exported Errors
var (
	ErrNotFound                       = fmt.Errorf("entity not found")
	ErrMissingManagingOrganization    = fmt.Errorf("missing managing organization")
	ErrMissingName                    = fmt.Errorf("missing name value")
	ErrMissingDescription             = fmt.Errorf("missing description value")
	ErrMalformedInputValue            = fmt.Errorf("malformed input value")
	ErrMissingOrganization            = fmt.Errorf("missing organization")
	ErrMissingProposition             = fmt.Errorf("missing proposition")
	ErrMissingGlobalReference         = fmt.Errorf("missing global reference")
	ErrNotImplementedByHSDP           = fmt.Errorf("method not implemented by HSDP")
	ErrCouldNoReadResourceAfterCreate = fmt.Errorf("could not read resource after create")
	ErrBaseIDMCannotBeEmpty           = fmt.Errorf("base IDM URL cannot be empty")
	ErrBaseIAMCannotBeEmpty           = fmt.Errorf("base IDM URL cannot be empty")
	ErrEmptyResults                   = fmt.Errorf("empty results")
	ErrOperationFailed                = fmt.Errorf("operation failed")
)
