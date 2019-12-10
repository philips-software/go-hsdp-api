package iam

import (
	"time"
)

type MFAPolicy struct {
	Schemas     []string          `json:"schemas" validate:"min=1"`
	ID          string            `json:"id,omitempty" validate:"omitempty,min=1,max=256"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Resource    MFAPolicyResource `json:"resource,omitempty"`
	ExternalID  string            `json:"externalId,omitempty"`
	Types       []string          `json:"types" validate:"min=1"`
	Active      *bool             `json:"active,omitempty"`
	CreatedBy   *struct {
		Value string `json:"value,omitempty"`
		Ref   string `json:"$ref,omitempty"`
	} `json:"createdBy,omitempty"`
	ModifiedBy *struct {
		Value string `json:"value,omitempty"`
		Ref   string `json:"$ref,omitempty"`
	} `json:"modifiedBy,omitempty"`
	Meta *MFAPolicyMeta `json:"meta,omitempty"`
}

type MFAPolicyResource struct {
	Type  string `json:"type" validate:"required"`
	Value string `json:"value" validate:"required"`
	Ref   string `json:"$ref,omitempty"`
}

type MFAPolicyMeta struct {
	ResourceType string     `json:"resourceType,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
	Location     string     `json:"location,omitempty"`
	Version      string     `json:"version,omitempty"`
}

func (p *MFAPolicy) SetActive(val bool) {
	p.Active = &val
}

func (p *MFAPolicy) SetResourceUser(uuid string) {
	p.Resource.Type = "User"
	p.Resource.Value = uuid
}

func (p *MFAPolicy) SetResourceOrganization(uuid string) {
	p.Resource.Type = "Organization"
	p.Resource.Value = uuid
}

func (p *MFAPolicy) SetType(val string) {
	p.Types = append([]string{}, val)
}
