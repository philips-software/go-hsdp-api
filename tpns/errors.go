package tpns

import (
	errors "golang.org/x/xerrors"
)

// Errors
var (
	ErrBaseTPNSCannotBeEmpty = errors.New("TPNS base URL cannot be empty")
)
