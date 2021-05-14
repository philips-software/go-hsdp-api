package notification

import (
	"errors"
)

// Errors
var (
	ErrNotificationURLCannotBeEmpty = errors.New("base Notification URL cannot be empty")
	ErrEmptyResult                  = errors.New("empty result")
	ErrNotificationForbidden        = errors.New("HTTP 403 Notification response")
	ErrBadRequest                   = errors.New("HTTP 400 Bad request")
	ErrNonHttp20xResponse           = errors.New("non HTTP 20x Notification response")
)
