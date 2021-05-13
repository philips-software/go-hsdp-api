package notification

import (
	"errors"
)

// Errors
var (
	ErrNotificationURLCannotBeEmpty = errors.New("base Notification URL cannot be empty")
	ErrEmptyResult                  = errors.New("empty result")
	ErrNotificationForbidden        = errors.New("HTTP 403 Notification response")
	ErrNonHttp20xResponse           = errors.New("non HTTP 20x Notification response")
)
