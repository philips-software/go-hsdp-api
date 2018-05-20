package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

// Proposition represents an IAM Proposition entity
type Proposition struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	OrganizationID    string `json:"organizationId"`
	GlobalReferenceID string `json:"globalReferenceId"`
}

func (p *Proposition) Validate() error {
	if p.Name == "" {
		return errMissingName
	}
	if p.OrganizationID == "" {
		return errMissingOrganization
	}
	if p.GlobalReferenceID == "" {
		return errMissingGlobalReference
	}
	return nil
}

func (a *Proposition) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0)
	a.ID, _ = r.Path("id").Data().(string)
	a.Name, _ = r.Path("name").Data().(string)
	a.Description, _ = r.Path("description").Data().(string)
	a.OrganizationID, _ = r.Path("organizationId").Data().(string)
	a.GlobalReferenceID, _ = r.Path("globalReferenceId").Data().(string)
	// TODO: Add new "meta" info as well
	return nil
}
