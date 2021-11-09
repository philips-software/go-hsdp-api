package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

func TestCreateApplication(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	appID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	propID := "3af7143e-de76-11e8-9681-6a0002b8cb70"
	description := "TESTPROP Application"
	globalReferenceID := "TESTAPPREF"

	muxMDM.HandleFunc("/connect/mdm/Application", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("_id") != appID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
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
			w.Header().Set("Location", "/connect/mdm/Application/"+appID)
			w.WriteHeader(http.StatusCreated)
		}
	})

	var app = mdm.Application{
		Description:       description,
		PropositionID:     propID,
		GlobalReferenceID: globalReferenceID,
		ApplicationGUID: mdm.Identifier{
			System: "foo",
			Value:  "bar",
		},
		DefaultGroupGUID: mdm.Identifier{
			System: "foo",
			Value:  "bar",
		},
	}

	// Test validation
	_, _, err := mdmClient.Applications.CreateApplication(app)
	if err == nil {
		t.Error("Expected validation error")
	}

	app.Name = "TESTAPP"
	createdApp, resp, err := mdmClient.Applications.CreateApplication(app)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.Equal(t, appID, createdApp.ID)
	assert.Equal(t, propID, createdApp.PropositionID)
	assert.Equal(t, globalReferenceID, createdApp.GlobalReferenceID)
}

func TestApplicationErrors(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	appID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	app, _, err := mdmClient.Applications.GetApplicationByID(appID)
	assert.NotNil(t, err)
	assert.Nil(t, app)
	app, _, err = mdmClient.Applications.GetApplicationByName("name")
	assert.NotNil(t, err)
	assert.Nil(t, app)
	apps, _, err := mdmClient.Applications.GetApplications(&mdm.GetApplicationsOptions{})
	assert.NotNil(t, err)
	assert.Nil(t, apps)
	app, _, err = mdmClient.Applications.CreateApplication(mdm.Application{})
	assert.NotNil(t, err)
	assert.Nil(t, app)
}
