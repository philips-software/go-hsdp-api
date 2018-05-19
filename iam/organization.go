package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

type Organization struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	DistinctName   string `json:"distinctName"`
	OrganizationID string `json:"organizationId"`
}

func (o *Organization) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0).Path("resource")
	o.OrganizationID, _ = r.Path("id").Data().(string)
	o.Name, _ = r.Path("name").Data().(string)
	o.Description, _ = r.Path("text").Data().(string)
	return nil
}
