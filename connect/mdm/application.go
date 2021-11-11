package mdm

import "time"

// Application represents an MDM Application entity
type Application struct {
	ID                string     `json:"id,omitempty"`
	ResourceType      string     `json:"resourceType"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	PropositionID     Reference  `json:"propositionId"`
	ApplicationGuid   Identifier `json:"applicationGuid"`
	GlobalReferenceID string     `json:"globalReferenceId"`
	DefaultGroupGuid  Identifier `json:"defaultGroupGuid"`
	Meta              *Meta      `json:"meta,omitempty"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	VersionID   string    `json:"versionId,omitempty"`
}

type Identifier struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value" validate:"required"`
}
