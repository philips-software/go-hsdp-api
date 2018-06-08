package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestGetOrganizationByID(t *testing.T) {
	teardown := setup()
	defer teardown()

	orgUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	orgName := "TestDevOrg"
	muxIDM.HandleFunc("/authorize/identity/Organization", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"resourceType": "bundle",
			"type": "searchset",
			"total": 1,
			"entry": [
				{
					"resource": {
						"id": "`+orgUUID+`",
						"resourceType": "Organization",
						"text": "TestDev",
						"name": "`+orgName+`"
					},
					"fullUrl": "https%3A%2F%2Fidm-staging.us-east.philips-healthsuite.com%2Fauthorize%2Fidentity%2FOrganization%3F_id%3D`+orgUUID+`"
				}
			]
		}`)
	})

	org, resp, err := client.Organizations.GetOrganizationByID(orgUUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if org.Name != orgName {
		t.Errorf("Expected Org name: %s, Got: %s", orgName, org.Name)
	}
	if org.OrganizationID != orgUUID {
		t.Errorf("Expected Org UUID: %s, Got: %s", orgUUID, org.OrganizationID)
	}
}
