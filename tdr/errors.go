package tdr

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrBaseTDRCannotBeEmpty           = errors.New("base TDR URL cannot be empty")
	ErrMissingTDRScopes               = errors.New("missing TDR scopes. Need 'tdr.contract' and 'tdr.dataitem'")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = fmt.Errorf("could not read resource after create")
)
