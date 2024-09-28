package cdr_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/fhir/go/fhirversion"
	"github.com/philips-software/go-hsdp-api/cdr"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/assert"
)

func TestR4PatchOperation(t *testing.T) {
	teardown := setup(t, fhirversion.R4)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
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
			body, err := io.ReadAll(r.Body)
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
	patched, resp, err := cdrClient.OperationsR4.Patch("Organization/"+orgID, []byte(`{"op": "replace","path": "/name","value": "Hospital2"}
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
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestR4PostOperation(t *testing.T) {
	teardown := setup(t, fhirversion.R4)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
		switch r.Method {
		case "POST":
			if !assert.Equal(t, "application/fhir+json;fhirVersion=4.0", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, cdr.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			body, err := io.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			contained, err := um.UnmarshalR4(body)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			onboardedOrg := contained.GetOrganization()
			jsonOrg, err := ma.MarshalResource(onboardedOrg)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(jsonOrg)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	patched, resp, err := cdrClient.OperationsR4.Post("Organization/"+orgID, []byte(`{
  "resourceType": "Organization",
  "id": "dae89cf0-888d-4a26-8c1d-578e97365efc",
  "meta": {
    "versionId": "4cbb8588-444a-11eb-917c-1f1d96935807",
    "lastUpdated": "2020-12-22T11:39:07.055441+00:00"
  },
  "identifier": [
    {
      "use": "usual",
      "system": "https://identity.philips-healthsuite.com/organization",
      "value": "dae89cf0-888d-4a26-8c1d-578e97365efc"
    }
  ],
  "name": "Hospital"
}`))

	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, patched) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}

func TestR4GetOperation(t *testing.T) {
	teardown := setup(t, fhirversion.R4)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
		switch r.Method {
		case "GET":
			if !assert.Equal(t, "application/fhir+json;fhirVersion=4.0", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, cdr.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
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
	retrieved, resp, err := cdrClient.OperationsR4.Get("Organization/" + orgID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, retrieved) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	org := retrieved.GetOrganization()
	assert.Equal(t, "Hospital2", org.Name.Value)
}

func TestR4DeleteOperation(t *testing.T) {
	teardown := setup(t, fhirversion.R4)
	defer teardown()

	orgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"

	muxCDR.HandleFunc("/store/fhir/"+cdrOrgID+"/Organization/"+orgID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
		switch r.Method {
		case "DELETE":
			if !assert.Equal(t, cdr.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ok, resp, err := cdrClient.OperationsR4.Delete("Organization/" + orgID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)
}
