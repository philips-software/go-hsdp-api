package cartel

import errors "golang.org/x/xerrors"

var (
	ErrMissingSecret      = errors.New("missing secret")
	ErrMissingToken       = errors.New("missing token")
	ErrMissingHost        = errors.New("missing host")
	ErrNotImplemented     = errors.New("not implemented")
	ErrNotFound           = errors.New("not found")
	ErrNonHttp20xResponse = errors.New("non http 20x response")
)
