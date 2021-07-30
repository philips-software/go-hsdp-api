package notification

import (
	"errors"
)

// Errors
var (
	ErrNotificationURLCannotBeEmpty = errors.New("base Notification URL cannot be empty")
	ErrEmptyResult                  = errors.New("empty result")
)
