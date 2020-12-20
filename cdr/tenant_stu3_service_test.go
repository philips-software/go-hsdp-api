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
			w.WriteHeader(http.StatusInternalServerError)
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
}
