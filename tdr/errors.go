package tdr

import (
	"errors"
)

// Errors
var (
	ErrBaseTDRCannotBeEmpty           = errors.New("base TDR URL cannot be empty")
	ErrMissingTDRScopes               = errors.New("missing TDR scopes. Need 'tdr.contract' and 'tdr.dataitem'")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrEmptyResults                   = errors.New("empty results")
)
