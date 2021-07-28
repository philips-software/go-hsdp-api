package cartel

import (
	"errors"
	"regexp"
)

var (
	ErrMissingSecret         = errors.New("missing cartel secret")
	ErrMissingToken          = errors.New("missing cartel token")
	ErrMissingHost           = errors.New("missing cartel host")
	ErrNotFound              = errors.New("not found")
	ErrHostnameAlreadyExists = errors.New("hostname already exists")
	ErrInvalidSubnetType     = errors.New("invalid subnet type, must be public or private")
)

var (
	existRegexErr = regexp.MustCompile(`^Host named [^\s]+ already exists!`)
)
