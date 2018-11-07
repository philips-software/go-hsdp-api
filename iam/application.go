package iam

import (
	"github.com/jeffail/gabs"
)

// Application represents an IAM Application entity
type Application struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name" validate:"required"`
	Description       string `json:"description"`
	PropositionID     string `json:"propositionId" validate:"required"`
	GlobalReferenceID string `json:"globalReferenceId" validate:"required"`
}

func (a *Application) parseFromBundle(bundle []byte) error {
	jsonParsed, err := gabs.ParseJSON(bundle)
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
