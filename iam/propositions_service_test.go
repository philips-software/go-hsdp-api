package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestCreateProposition(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	propID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"
	orgID := "3af7143e-de76-11e8-9681-6a0002b8cb70"
	description := "TESTPROP Proposition"

	muxIDM.HandleFunc("/authorize/identity/Proposition", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("_id") != propID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
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
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", "/authorize/identity/Proposition/"+propID)
			w.WriteHeader(http.StatusCreated)
		}
	})

	var prop = Proposition{
		Name:              "TESTPROP",
		Description:       description,
		OrganizationID:    orgID,
		GlobalReferenceID: "TESTPROPREF",
	}
	createdProp, resp, err := client.Propositions.CreateProposition(prop)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP created")
	}
	if createdProp.ID != propID {
		t.Errorf("Unexpected ID")
	}
}

func TestGetPropositions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "168cdeae-2539-45b0-b18c-89ae32f1ea15"
	propID := "10dc5e2f-3940-4cd8-b0ef-297e12ad2f3c"

	muxIDM.HandleFunc("/authorize/identity/Proposition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
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

	props, resp, err := client.Propositions.GetPropositions(&GetPropositionsOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if len(*props) != 1 {
		t.Errorf("Expected 1 proposition, Got: %d", len(*props))
	}
}
