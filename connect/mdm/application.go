package mdm

import "time"

// Application represents an MDM Application entity
type Application struct {
	ResourceType      string     `json:"resourceType"`
	ID                string     `json:"id,omitempty"`
	Name              string     `json:"name" validate:"required"`
	Description       string     `json:"description"`
	PropositionID     string     `json:"propositionId" validate:"required"`
	GlobalReferenceID string     `json:"globalReferenceId" validate:"required"`
	ApplicationGUID   Identifier `json:"applicationGuid"`
	DefaultGroupGUID  Identifier `json:"defaultGroupGuid"`
	Meta              *Meta      `json:"meta,omitempty"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	VersionID   string    `json:"versionId,omitempty"`
}

type Identifier struct {
	System string `json:"system"`
	Value  string `json:"value" validate:"required"`
}
