package discovery

import (
	"errors"
)

var (
	ErrBaseURLCannotBeEmpty = errors.New("base URL cannot be empty")
	ErrEmptyResult          = errors.New("empty result")
	ErrInvalidEndpointURL   = errors.New("invalid endpoint URL")
)
