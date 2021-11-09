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

	muxMDM.HandleFunc("/connect/mdm/Proposition", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("_id") != propID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
                                     "total": 1,
                                     "entry": [
                                       {
                                         "name": "TESTPROP",
                                         "description": "`+description+`",
                                         "organizationId": "`+orgID+`",
                                         "globalReferenceId": "TESTPROP",
                                         "id": "`+propID+`",
                                         "meta": {
                                           "versionId": "0",
                                           "lastModified": "2018-11-02T05:48:410.042Z"
                                         }
                                       }
                                     ]
                                   }`)
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", "/connect/mdm/Proposition/"+propID)
			w.WriteHeader(http.StatusCreated)
		}
	})

	var prop = mdm.Proposition{
		Name:              "TESTPROP",
		Description:       description,
		OrganizationID:    orgID,
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
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestGetPropositions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "168cdeae-2539-45b0-b18c-89ae32f1ea15"
	propID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"

	muxMDM.HandleFunc("/connect/mdm/Proposition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"total": 1,
			"entry": [
				{
					"name": "TESTPROP",
					"description": "TEST Proposition",
					"organizationId": "`+orgID+`",
					"globalReferenceId": "testprop-1",
					"id": "`+propID+`",
					"meta": {
						"versionId": "0",
						"lastModified": "2018-06-28T08:41:895.010Z"
					}
				}
			]
		}`)
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
