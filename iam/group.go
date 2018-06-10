package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

// Group represents an IAM group resource
type Group struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	ManagingOrganization string `json:"managingOrganization,omitempty"`
}

func (g *Group) validate() error {
	if g.ManagingOrganization == "" {
		return errMissingManagingOrganization
	}
	if g.Name == "" {
		return errMissingName
	}
	return nil
}

func (g *Group) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0).Path("resource")
	g.ID = r.Path("_id").Data().(string)
	g.ManagingOrganization, _ = r.Path("orgId").Data().(string)
	g.Name, _ = r.Path("groupName").Data().(string)
	g.Description, _ = r.Path("groupDescription").Data().(string)
	return nil
}
