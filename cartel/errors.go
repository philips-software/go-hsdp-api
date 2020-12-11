package cartel

import "errors"

var (
	ErrMissingSecret      = errors.New("missing cartel secret")
	ErrMissingToken       = errors.New("missing cartel token")
	ErrMissingHost        = errors.New("missing cartel host")
	ErrNotImplemented     = errors.New("not implemented by cartel client")
	ErrNotFound           = errors.New("not found")
	ErrNonHttp20xResponse = errors.New("non http 20x cartel response")
	ErrInvalidSubnetType  = errors.New("invalid subnet type, must be public or private")
)
