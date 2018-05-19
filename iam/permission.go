package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Type        string `json:"type"`
}

func (p *Permission) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0)
	p.ID = r.Path("id").Data().(string)
	p.Category, _ = r.Path("category").Data().(string)
	p.Name, _ = r.Path("name").Data().(string)
	p.Description, _ = r.Path("description").Data().(string)
	p.Type, _ = r.Path("type").Data().(string)
	return nil
}
