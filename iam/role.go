package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

type Role struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	ManagingOrganization string `json:"managingOrganization"`
}

func (p *Role) parseFromBundle(v interface{}) error {
	m, err := json.Marshal(v)
	if err != nil {
		return err
	}
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0)
	p.ID = r.Path("id").Data().(string)
	p.ManagingOrganization, _ = r.Path("managingOrganization").Data().(string)
	p.Name, _ = r.Path("name").Data().(string)
	p.Description, _ = r.Path("description").Data().(string)
	return nil
}
