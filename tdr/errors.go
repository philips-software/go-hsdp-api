package tdr

import "fmt"

var (
	errBaseTDRCannotBeEmpty = fmt.Errorf("base TDR URL cannot be empty")
	errMissingTDRScopes     = fmt.Errorf("missing TDR scopes. Need 'tdr.contract' and 'tdr.dataitem'")
)
