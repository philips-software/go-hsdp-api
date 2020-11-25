package has

import (
	"errors"
)

// Errors
var (
	ErrBaseHASCannotBeEmpty           = errors.New("base HAS URL cannot be empty")
	ErrMissingHASPermissions          = errors.New("missing HAS permissions. Need 'HAS_RESOURCE.ALL' and 'HAS_SESSION.ALL'")
	ErrEmptyResult                    = errors.New("empty result")
	ErrCouldNoReadResourceAfterCreate = errors.New("could not read resource after create")
	ErrEmptyResults                   = errors.New("empty results")
)
