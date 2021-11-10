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
	propID := "Proposition/d9001b60-e4dc-420a-8176-333186898030"
	description := "TESTPROP Application"
	globalReferenceID := "TESTAPPREF"
	bundleResponse := `{
  "meta": {
    "lastUpdated": "2021-11-10T00:10:09.479525+00:00"
  },
  "id": "51951097-eecb-4b03-9e1b-68d1dc437cb7",
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "search": {
        "mode": "match"
      },
      "resource": {
        "meta": {
          "lastUpdated": "2021-11-09T23:36:40.988817+00:00",
          "versionId": "77d25cb2-d035-4dab-b20d-333133a3995d"
        },
        "id": "` + appID + `",
        "resourceType": "Application",
        "name": "First",
        "description": "First app",
        "propositionId": {
          "reference": "` + propID + `"
        },
        "applicationGuid": {
          "system": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/identity/",
          "value": "ca277a24-22a4-4dd0-b8d5-bc961b74635a"
        },
        "globalReferenceId": "` + globalReferenceID + `",
        "defaultGroupGuid": {
          "system": "https://idm-client-test.us-east.philips-healthsuite.com",
          "value": "060e2ec0-8775-41db-8372-007de2b7dbef"
        }
      },
      "fullUrl": "Application/` + appID + `"
    }
  ],
  "link": [
    {
      "url": "Application?name=First&propositionId=d9001b60-e4dc-420a-8176-333186898030",
      "relation": "self"
    },
    {
      "url": "Application?name=First&propositionId=d9001b60-e4dc-420a-8176-333186898030&_page=1",
      "relation": "first"
    }
  ],
  "pageTotal": 1
}`

	muxMDM.HandleFunc("/connect/mdm/Application", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("_id") != appID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bundleResponse)
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", "/connect/mdm/Application/"+appID)
			w.WriteHeader(http.StatusCreated)
		}
	})

	var app = mdm.Application{
		Description: description,
		PropositionID: mdm.Reference{
			Reference: propID,
		},
		GlobalReferenceID: globalReferenceID,
		ApplicationGuid: mdm.Identifier{
			System: "foo",
			Value:  "bar",
		},
		DefaultGroupGuid: mdm.Identifier{
			System: "foo",
			Value:  "bar",
		},
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
	assert.Equal(t, propID, createdApp.PropositionID.Reference)
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
