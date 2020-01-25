package iam

import (
	errors "golang.org/x/xerrors"
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
	ErrNotImplementedByHSDP           = errors.New("method not implemented by HSDP")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrBaseIDMCannotBeEmpty           = errors.New("base IDM URL cannot be empty")
	ErrBaseIAMCannotBeEmpty           = errors.New("base IDM URL cannot be empty")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrMissingEtagInformation         = errors.New("missing etag information")
	ErrMissingRefreshToken            = errors.New("missing refresh token")
)
