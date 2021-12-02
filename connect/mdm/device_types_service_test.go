package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

func TestDeviceTypesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	name := "TestDeviceType"
	description := "Device Type description"
	createdResource := `{
  "meta": {
    "lastUpdated": "2021-11-11T11:07:42.351502+00:00",
    "versionId": "449be385-9799-49c8-893e-c53d51f1e6ce"
  },
  "id": "` + id + `",
  "resourceType": "DeviceType",
  "name": "` + name + `",
  "description": "` + description + `",
  "deviceGroupId": {
    "reference": "DeviceGroup/73c8df4a-703f-42ba-bdf2-40577a22e690"
  },
  "ctn": "WATCH1"
}`
	muxMDM.HandleFunc("/connect/mdm/DeviceType", func(w http.ResponseWriter, r *http.Request) {
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
	muxMDM.HandleFunc("/connect/mdm/DeviceType/"+id, func(w http.ResponseWriter, r *http.Request) {
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

	var c mdm.DeviceType
	c.Name = name
	c.Description = description
	c.CTN = "WATCH1"

	created, resp, err := mdmClient.DeviceTypes.Create(c)
	if !assert.Nilf(t, err, "unexpected error: %v", err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, name, created.Name)

	created, resp, err = mdmClient.DeviceTypes.GetByID(created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdResource) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, id, created.ID)

	ok, resp, err := mdmClient.DeviceTypes.Delete(*created)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdResource)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
