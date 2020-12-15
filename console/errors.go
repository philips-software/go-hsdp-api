package console

import "errors"

// Exported Errors
var (
	ErrNotFound                       = errors.New("entity not found")
	ErrMissingName                    = errors.New("missing name value")
	ErrMissingDescription             = errors.New("missing description value")
	ErrMalformedInputValue            = errors.New("malformed input value")
	ErrNotImplementedByHSDP           = errors.New("method not implemented by HSDP")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrConsoleURLCannotBeEmpty        = errors.New("console base URL cannot be empty")
	ErrUAAURLCannotBeEmpty            = errors.New("UAA URL cannot be empty")
	ErrEmptyResults                   = errors.New("empty results")
	ErrOperationFailed                = errors.New("operation failed")
	ErrMissingEtagInformation         = errors.New("missing etag information")
	ErrMissingRefreshToken            = errors.New("missing refresh token")
	ErrNotAuthorized                  = errors.New("not authorized")
	ErrNonHttp20xResponse             = errors.New("non http 20x console response")
)
