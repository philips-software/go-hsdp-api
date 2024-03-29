package iam

import (
	"time"
)

type OrgAddress struct {
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Country       string `json:"country,omitempty"`
}

type Meta struct {
	ResourceType string     `json:"resourceType,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
	LastUpdated  *time.Time `json:"lastUpdated,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
	UpdatedBy    string     `json:"updatedBy,omitempty"`
	CreatedBy    string     `json:"createdBy,omitempty"`
	Location     string     `json:"location,omitempty"`
	Version      string     `json:"version,omitempty"`   // No
	VersionID    string     `json:"versionId,omitempty"` // Consistency
}

type Attribute struct {
	Value   string `json:"value,omitempty"`
	Ref     string `json:"$ref,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// Organization represents a IAM Organization resource
type Organization struct {
	Schemas           []string    `json:"schemas"`
	ID                string      `json:"id"`
	ExternalID        string      `json:"externalId,omitempty"`
	Name              string      `json:"name"`
	DisplayName       string      `json:"displayName,omitempty"`
	Description       string      `json:"description,omitempty"`
	Parent            Attribute   `json:"parent,omitempty"`
	Type              string      `json:"type,omitempty"`
	Active            bool        `json:"active,omitempty"`
	InheritProperties bool        `json:"inheritProperties,omitempty"`
	Address           OrgAddress  `json:"address,omitempty"`
	Owners            []Attribute `json:"owners,omitempty"`
	CreatedBy         *Attribute  `json:"createdBy,omitempty"`
	ModifiedBy        *Attribute  `json:"modifiedBy,omitempty"`
	Meta              *Meta       `json:"meta,omitempty"`
}

type OrganizationStatus struct {
	Schemas        []string `json:"schemas"`
	ID             string   `json:"id"`
	Status         string   `json:"status"`
	TotalResources int      `json:"totalResources"`
	Meta           *Meta    `json:"meta"`
}
