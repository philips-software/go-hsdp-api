package ai

import (
	"errors"
)

// Errors
var (
	ErrAnalyzeURLCannotBeEmpty = errors.New("base Inference URL cannot be empty")
	ErrEmptyResult             = errors.New("empty result")
	ErrInvalidEndpointURL      = errors.New("invalid endpoint URL")
)
