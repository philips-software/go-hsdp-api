package dicom

import (
	"errors"
)

// Errors
var (
	ErrDICOMURLCannotBeEmpty = errors.New("base DICOM URL cannot be empty")
	ErrEmptyResult           = errors.New("empty result")
	ErrNonHttp20xResponse    = errors.New("non HTTP 20x DICOM response")
)
