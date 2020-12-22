package cdr_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/cdr"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/assert"
)

func TestPatchOperation(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		switch r.Method {
		case "PATCH":
			if !assert.Equal(t, "application/json-patch+json", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, cdr.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			body, err := ioutil.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = jsonpatch.MergePatch([]byte(`{}`), body)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
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
  "name": "Hospital2"
}
`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	patched, resp, err := cdrClient.OperationsSTU3.Patch("Organization/"+orgID, []byte(`{"op": "replace","path": "/name","value": "Hospital2"}
`))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, patched) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
