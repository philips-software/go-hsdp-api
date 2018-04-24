package api

import "fmt"

var (
	errNotFound                    = fmt.Errorf("Entity not found")
	errMissingManagingOrganization = fmt.Errorf("Missing managing organization")
	errMissingName                 = fmt.Errorf("Missing name value")
	errMissingDescription          = fmt.Errorf("Missing description value")
	errMalformedInputValue         = fmt.Errorf("Malformed input value")
)
