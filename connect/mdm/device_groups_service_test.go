package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

// TODO replace with DeviceGroup capture
func TestDeviceGroupCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	name := "TestService"
	description := "Service description"
	organizationID := "c3fe79e6-13c2-48c1-adfa-826a01d4b31c"
	createdResource := `{
  "meta": {
    "lastUpdated": "2021-11-10T18:47:24.059503+00:00",
    "versionId": "ae6fabc0-3a65-4bb2-a08f-64d4e6453f0d"
  },
  "id": "` + id + `",
  "resourceType": "ServiceAction",
  "name": "` + name + `",
  "description": "` + description + `",
  "ServiceActionId": {
    "reference": "ServiceAction/7b26ddb7-910b-4faf-b122-e1fd27356b14"
  },
  "organizationGuid": {
    "system": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/identity/Organization",
    "value": "` + organizationID + `"
  }
}`
	muxMDM.HandleFunc("/connect/mdm/ServiceAction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "pageTotal": 0,
  "link": [
    {
      "relation": "string",
      "url": "string"
    }
  ],
  "entry": [
    {
      "fullUrl": "string",
      "resource": `+createdResource+`
    }
  ]
}`)
		}
	})
	muxMDM.HandleFunc("/connect/mdm/ServiceAction/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var c mdm.ServiceAction
	c.Name = name
	c.Description = description
	c.OrganizationGuid = &mdm.Identifier{
		Value: organizationID,
	}

	created, resp, err := mdmClient.ServiceActions.Create(c)
	if !assert.Nilf(t, err, "unexpected error: %v", err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, name, created.Name)

	created, resp, err = mdmClient.ServiceActions.GetByID(created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdResource) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, id, created.ID)

	ok, resp, err := mdmClient.ServiceActions.Delete(*created)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdResource)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}
