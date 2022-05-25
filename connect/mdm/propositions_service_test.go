package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

func TestCreateProposition(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	propID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	orgID := "3af7143e-de76-11e8-9681-6a0002b8cb70"
	description := "TESTPROP Proposition"
	bundleResponse := `{
  "meta": {
    "lastUpdated": "2021-11-09T23:49:05.035529+00:00"
  },
  "id": "e1f40cd1-f821-4f25-a2d2-a13abf614a1f",
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "search": {
        "mode": "match"
      },
      "resource": {
        "meta": {
          "lastUpdated": "2021-11-09T23:37:13.643501+00:00",
          "versionId": "db393c1f-c8e8-4609-b8ae-143584010326"
        },
        "id": "` + propID + `",
        "resourceType": "Proposition",
        "name": "First",
        "description": "Description here",
        "organizationGuid": {
          "system": "https://iam-client-test.us-east.philips-healthsuite.com",
          "value": "` + orgID + `"
        },
        "propositionGuid": {
          "system": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/identity/",
          "value": "` + propID + `"
        },
        "globalReferenceId": "be5ea8a2-e8ad-483b-a6d2-77e38a6c25b9",
        "defaultCustomerOrganizationGuid": {
          "system": "https://iam-client-test.us-east.philips-healthsuite.com",
          "value": "` + orgID + `"
        },
        "status": "ACTIVE",
        "validationEnabled": false,
        "notificationEnabled": false
      },
      "fullUrl": "Proposition/d9001b60-e4dc-420a-8176-333186898030"
    }
  ],
  "link": [
    {
      "url": "Proposition?name=First&organizationGuid=` + orgID + `",
      "relation": "self"
    },
    {
      "url": "Proposition?name=First&organizationGuid=` + orgID + `&_page=1",
      "relation": "first"
    }
  ],
  "pageTotal": 1
}`

	muxMDM.HandleFunc("/connect/mdm/Proposition", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("_id") != propID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bundleResponse)
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
  "meta": {
    "lastUpdated": "2021-11-09T23:37:13.643501+00:00",
    "versionId": "db393c1f-c8e8-4609-b8ae-143584010326"
  },
  "id": "`+propID+`",
  "resourceType": "Proposition",
  "name": "First",
  "description": "Description here",
  "organizationGuid": {
    "system": "https://iam-client-test.us-east.philips-healthsuite.com",
    "value": "`+orgID+`"
  },
  "propositionGuid": {
    "system": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/identity/",
    "value": "64e403e6-d215-457a-bf12-2a4f49038208"
  },
  "globalReferenceId": "be5ea8a2-e8ad-483b-a6d2-77e38a6c25b9",
  "defaultCustomerOrganizationGuid": {
    "system": "https://iam-client-test.us-east.philips-healthsuite.com",
    "value": "`+orgID+`"
  },
  "status": "ACTIVE",
  "validationEnabled": false,
  "notificationEnabled": false
}`)
		}
	})

	var prop = mdm.Proposition{
		Name:              "TESTPROP",
		Description:       description,
		OrganizationGuid:  mdm.Identifier{Value: orgID},
		GlobalReferenceID: "TESTPROPREF",
	}
	createdProp, resp, err := mdmClient.Propositions.CreateProposition(prop)
	if err != nil {
		t.Fatal(err)
	}
	if ok := assert.NotNil(t, createdProp); ok {
		assert.Equal(t, propID, createdProp.ID)
	}
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}
}

func TestGetPropositions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "168cdeae-2539-45b0-b18c-89ae32f1ea15"
	propID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	bundleResponse := `{
  "meta": {
    "lastUpdated": "2021-11-09T23:49:05.035529+00:00"
  },
  "id": "e1f40cd1-f821-4f25-a2d2-a13abf614a1f",
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "search": {
        "mode": "match"
      },
      "resource": {
        "meta": {
          "lastUpdated": "2021-11-09T23:37:13.643501+00:00",
          "versionId": "db393c1f-c8e8-4609-b8ae-143584010326"
        },
        "id": "` + propID + `",
        "resourceType": "Proposition",
        "name": "First",
        "description": "Description here",
        "organizationGuid": {
          "system": "https://iam-client-test.us-east.philips-healthsuite.com",
          "value": "` + orgID + `"
        },
        "propositionGuid": {
          "system": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/identity/",
          "value": "` + propID + `"
        },
        "globalReferenceId": "be5ea8a2-e8ad-483b-a6d2-77e38a6c25b9",
        "defaultCustomerOrganizationGuid": {
          "system": "https://iam-client-test.us-east.philips-healthsuite.com",
          "value": "` + orgID + `"
        },
        "status": "ACTIVE",
        "validationEnabled": false,
        "notificationEnabled": false
      },
      "fullUrl": "Proposition/d9001b60-e4dc-420a-8176-333186898030"
    }
  ],
  "link": [
    {
      "url": "Proposition?name=First&organizationGuid=dae89cf0-888d-4a26-8c1d-578e97365efc",
      "relation": "self"
    },
    {
      "url": "Proposition?name=First&organizationGuid=dae89cf0-888d-4a26-8c1d-578e97365efc&_page=1",
      "relation": "first"
    }
  ],
  "pageTotal": 1
}`

	muxMDM.HandleFunc("/connect/mdm/Proposition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, bundleResponse)
	})

	props, resp, err := mdmClient.Propositions.GetPropositions(&mdm.GetPropositionsOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if ok := assert.NotNil(t, props); ok {
		assert.Equal(t, 1, len(*props))
	}
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}
