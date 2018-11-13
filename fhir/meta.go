package fhir

type Meta struct {
	Element
	VersionId   string        `json:"versionId,omitempty"`
	LastUpdated *FHIRDateTime `json:"lastUpdated,omitempty"`
	Profile     []string      `json:"profile,omitempty"`
	Security    []Coding      `json:"security,omitempty"`
	Tag         []Coding      `json:"tag,omitempty"`
}
