package cce

import (
	errors "golang.org/x/xerrors"
)

// Exported Errors
var (
	ErrNotFound                       = errors.New("entity not found")
	ErrMissingName                    = errors.New("missing name value")
	ErrMalformedInputValue            = errors.New("malformed input value")
	ErrMissingOrganization            = errors.New("missing organization")
	ErrMissingGlobalReference         = errors.New("missing global reference")
	ErrNotImplementedByCCE            = errors.New("method not implemented by CCE")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrBaseURLCannotBeEmpty           = errors.New("base URL cannot be empty")
	ErrBaseIAMCannotBeEmpty           = errors.New("base IAM URL cannot be empty")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrMissingEtagInformation         = errors.New("missing etag information")
	ErrMissingRefreshToken            = errors.New("missing refresh token")
	ErrNotAuthorized                  = errors.New("not authorized")
	ErrNoValidSignerAvailable         = errors.New("no valid HSDP signer available")
	ErrDiscoveryFailed                = errors.New("discovery of CCE endpoints failed")
)

type UserError struct {
	User string
	Err  error
}

func (e *UserError) Error() string { return "user: " + e.User }

func (e *UserError) Unwrap() error { return e.Err }
