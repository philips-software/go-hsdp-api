package pki

import (
	errors "golang.org/x/xerrors"
)

// Errors
var (
	ErrBaseHASCannotBeEmpty           = errors.New("base PKI URL cannot be empty")
	ErrMissingPKIPermissions          = errors.New("missing PKI permissions")
	ErrMissingIAMOrganization         = errors.New("missing IAM organization")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrCertificateExpected            = errors.New("certificate expected")
	ErrCRLExpected                    = errors.New("certificate revocation list expected")
	ErrCFClientNotConfigured          = errors.New("CF client not configured")
	ErrCFInvalidToken                 = errors.New("invalid CF token")
)
