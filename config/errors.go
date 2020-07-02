package config

import (
	errors "golang.org/x/xerrors"
)

var (
	ErrMissingConfig                     = errors.New("missing config")
	ErrNotFound                          = errors.New("not found")
	ErrUnreachableOrOutdatedConfigSource = errors.New("unreachable or outdated config source")
)
