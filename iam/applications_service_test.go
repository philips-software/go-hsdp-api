package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestCreateApplication(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	appID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	propID := "3af7143e-de76-11e8-9681-6a0002b8cb70"
	description := "TESTPROP Application"
	globalReferenceID := "TESTAPPREF"

	muxIDM.HandleFunc("/authorize/identity/Application", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("_id") != appID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
                                     "total": 1,
                                     "entry": [
                                       {
                                         "name": "TESTAPP",
                                         "description": "`+description+`",
                                         "propositionId": "`+propID+`",
                                         "globalReferenceId": "`+globalReferenceID+`",
                                         "id": "`+appID+`",
                                         "meta": {
                                           "versionId": "0",
                                           "lastModified": "2018-11-02T05:48:410.042Z"
                                         }
                                       }
                                     ]
                                   }`)
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", "/authorize/identity/Application/"+appID)
			w.WriteHeader(http.StatusCreated)
		}
	})

	var app = Application{
		Description:       description,
		PropositionID:     propID,
		GlobalReferenceID: globalReferenceID,
	}

	// Test validation
	createdApp, resp, err := client.Applications.CreateApplication(app)
	if err == nil {
		t.Error("Expected validation error")
	}

	app.Name = "TESTAPP"
	createdApp, resp, err = client.Applications.CreateApplication(app)
	if err != nil {
		t.Fatal(err)
	}
	if createdApp.ID != appID {
		t.Error("Expected ID to be set")
	}
	if createdApp.PropositionID != propID {
		t.Error("Expected propositionID to be set")
	}
	if createdApp.GlobalReferenceID != globalReferenceID {
		t.Error("Expected GlobalReferenceID to be set")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP created")
	}
	if createdApp.ID != appID {
		t.Errorf("Unexpected ID")
	}
}
