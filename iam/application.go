package iam

import (
	"encoding/json"

	"github.com/jeffail/gabs"
)

// Application represents an IAM Application entity
type Application struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DistinctName      string `json:"distinctName"`
	PropositionID     string `json:"propositionId"`
	GlobalReferenceID string `json:"globalReferenceId"`
}

func (a *Application) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0)
	a.ID, _ = r.Path("id").Data().(string)
	a.Name, _ = r.Path("name").Data().(string)
	a.Description, _ = r.Path("description").Data().(string)
	a.PropositionID, _ = r.Path("propositionId").Data().(string)
	a.GlobalReferenceID, _ = r.Path("globalReferenceId").Data().(string)
	// TODO: Add new "meta" info as well
	return nil
}
