package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	parentOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	orgName := "TestDevOrg"
	orgDescription := "Some description"
	orgID := "af5eee7c-0203-4d5a-a021-414531d0f451"
	muxIDM.HandleFunc("/security/organizations/"+parentOrgID+"/childorganizations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"exchange": {
				"name": "`+orgName+`",
				"description": "`+orgDescription+`",
				"distinctName": "ou=`+orgName+`,ou=ToplevelORG,dc=foo-bar,dc=com",
				"organizationId": "`+orgID+`"
			},
			"responseCode": "200",
			"responseMessage": "Success"
		}`)
	})

	org, resp, err := client.Organizations.CreateOrganization(parentOrgID, orgName, orgDescription)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if org.Name != orgName {
		t.Errorf("Expected Org name: %s, Got: %s", orgName, org.Name)
	}
	if org.OrganizationID != orgID {
		t.Errorf("Expected Org UUID: %s, Got: %s", orgID, org.OrganizationID)
	}
}

func TestGetOrganizationByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	orgName := "TestDevOrg"
	muxIDM.HandleFunc("/authorize/identity/Organization", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		if ok, err := signerHSDP.ValidateRequest(r); !ok {
			t.Fatalf("Expected valid HSDP signature. Error: %v", err)
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
					"fullUrl": "https%3A%2F%2Fidm-something.foo-bar.com%2Fauthorize%2Fidentity%2FOrganization%3F_id%3D`+orgUUID+`"
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

func TestUpdateOrganization(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	orgName := "TestDevOrg"
	description := "New description"
	muxIDM.HandleFunc("/security/organizations/"+orgUUID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"exchange": {
				"name": "TEST2",
				"description": "`+description+`",
				"organizationId": "`+orgUUID+`"
			},
			"responseCode": "200",
			"responseMessage": "Success"
		}`)
	})
	var org Organization
	org.OrganizationID = orgUUID
	org.Description = description
	org.Name = orgName

	_, resp, err := client.Organizations.UpdateOrganization(org)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
}
