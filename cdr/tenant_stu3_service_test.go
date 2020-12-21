package cdr_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"

	"github.com/stretchr/testify/assert"
)

func TestTenantService(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json+fhir")
		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, string(body))
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Organization",
  "id": "`+orgID+`",
  "meta": {
    "versionId": "6dfa7cc8-2000-11ea-91df-bb500f85c5e2",
    "lastUpdated": "2019-12-16T12:34:40.544022+00:00"
  },
  "identifier": [
    {
      "use": "usual",
      "system": "https://identity.philips-healthsuite.com/organization",
      "value": "`+orgID+`"
    }
  ],
  "active": true,
  "name": "Hospital"
}
`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	org, err := stu3.NewOrganization(timeZone, orgID, "Hospital")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, org) {
		return
	}
	newOrg, resp, err := cdrClient.TenantSTU3.Onboard(org)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, newOrg) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	foundOrg, resp, err := cdrClient.TenantSTU3.GetOrganizationByID(orgID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, foundOrg) {
		return
	}
	assert.Equal(t, "Hospital", foundOrg.Name.Value)
}
