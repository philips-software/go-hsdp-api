package inference

import (
	"errors"
)

// Errors
var (
	ErrInferenceURLCannotBeEmpty = errors.New("base Inference URL cannot be empty")
	ErrEmptyResult               = errors.New("empty result")
)
