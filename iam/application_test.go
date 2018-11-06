package iam

import (
	"testing"
)

func TestApplication(t *testing.T) {
	var app Application

	// Missing Name
	if err := app.Validate(); err != ErrMissingName {
		t.Errorf("Unexpected error: %v", err)
	}
	app.Name = "Name"

	// Missing PropositionID
	if err := app.Validate(); err != ErrMissingProposition {
		t.Errorf("Unexpected error: %v", err)
	}
	app.PropositionID = "uuid"

	// Missing GlobalReferenceID
	if err := app.Validate(); err != ErrMissingGlobalReference {
		t.Errorf("Unexpected error: %v", err)
	}
	app.GlobalReferenceID = "ref"

	err := app.Validate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseApplicationFromBundle(t *testing.T) {
	var app Application

	body := []byte(`{
       "total": 1,
       "entry": [
         {
           "name": "FOO",
           "description": "FOO - Development",
           "globalReferenceId": "6b8ef89a-de86-11e8-94fc-6a0002b8cb70",
           "propositionId": "5c168ce8-de86-11e8-b39d-6a0002b8cb70",
           "id": "65880950-de86-11e8-b804-6a0002b8cb70",
           "meta": {
             "versionId": "0",
             "lastModified": "2018-07-26T16:21:202.052Z"
           }
         }
       ]
     }`)

	err := app.parseFromBundle(body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if app.Name != "FOO" {
		t.Errorf("Unexpected name: %s, expected: FOO", app.Name)
	}
}
