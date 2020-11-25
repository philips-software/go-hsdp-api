package iam

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
	ErrNotImplementedByHSDP           = errors.New("method not implemented by HSDP")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrBaseIDMCannotBeEmpty           = errors.New("base IDM URL cannot be empty")
	ErrBaseIAMCannotBeEmpty           = errors.New("base IDM URL cannot be empty")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrMissingEtagInformation         = errors.New("missing etag information")
	ErrMissingRefreshToken            = errors.New("missing refresh token")
	ErrNotAuthorized                  = errors.New("not authorized")
	ErrNoValidSignerAvailable         = errors.New("no valid HSDP signer available")
)

type UserError struct {
	User string
	Err  error
}

func (e *UserError) Error() string { return "user: " + e.User }

func (e *UserError) Unwrap() error { return e.Err }
