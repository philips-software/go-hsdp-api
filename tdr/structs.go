package tdr

import "fmt"

// DataType describes the system and code of a resource
type DataType struct {
	System string `json:"system,omitempty"`
	Code   string `json:"code,omitempty"`
}

// Meta contains versioning info about the resource
type Meta struct {
	LastUpdated string `json:"lastUpdated,omitempty"`
	VersionID   string `json:"versionId,omitempty"`
}

// String pretty prints a DataType
func (d DataType) String() string {
	return fmt.Sprintf("tdr.DataType:System=%s,Code=%s", d.System, d.Code)
}
