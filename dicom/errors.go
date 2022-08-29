package dicom

import (
	"errors"
)

// Errors
var (
	ErrDICOMURLCannotBeEmpty = errors.New("base DICOM URL cannot be empty")
	ErrEmptyResult           = errors.New("empty result")
	ErrNotFound              = errors.New("not found")
)
