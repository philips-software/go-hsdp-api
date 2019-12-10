package iam

import (
	"encoding/json"
	"testing"
)

func TestApplication(t *testing.T) {
	var app Application

	body := []byte(`{
           "name": "FOO",
           "description": "FOO - Development",
           "globalReferenceId": "6b8ef89a-de86-11e8-94fc-6a0002b8cb70",
           "propositionId": "5c168ce8-de86-11e8-b39d-6a0002b8cb70",
           "id": "65880950-de86-11e8-b804-6a0002b8cb70",
           "meta": {
             "versionId": "0",
             "lastModified": "2018-07-26T16:21:202.052Z"
	    }
	}`)

	err := json.Unmarshal(body, &app)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if app.Name != "FOO" {
		t.Errorf("Unexpected name: %s, expected: FOO", app.Name)
	}
}
