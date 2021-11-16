package mdm

import "time"

// Application represents an MDM Application entity
type Application struct {
	ResourceType      string      `json:"resourceType" validate:"required"`
	ID                string      `json:"id,omitempty"`
	Name              string      `json:"name" validate:"required"`
	Description       string      `json:"description"`
	PropositionID     Reference   `json:"propositionId"`
	GlobalReferenceID string      `json:"globalReferenceId"`
	ApplicationGuid   *Identifier `json:"applicationGuid,omitempty"`
	DefaultGroupGuid  *Identifier `json:"defaultGroupGuid,omitempty"`
	Meta              *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	VersionID   string    `json:"versionId,omitempty"`
}

type Identifier struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value" validate:"required"`
}
