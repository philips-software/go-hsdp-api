package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

func TestServiceActionsCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	name := "TestService"
	description := "Service description"
	organizationID := "c3fe79e6-13c2-48c1-adfa-826a01d4b31c"
	createdService := `{
  "resourceType": "StandardService",
  "id": "` + id + `",
  "name": "` + name + `",
  "description": "` + description + `",
  "trusted": true,
  "tags": [
    "string"
  ],
  "serviceUrls": [
    {
      "url": "string",
      "sortOrder": 0,
      "authenticationMethodId": {
        "reference": "string"
      }
    }
  ],
  "meta": {
    "lastUpdated": "2021-11-09T22:15:35.155Z",
    "versionId": "string"
  },
  "organizationGuid": {
    "system": "string",
    "value": "` + organizationID + `"
  }
}
`
	muxMDM.HandleFunc("/connect/mdm/StandardService", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, createdService)
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
      "resource": {
        "resourceType": "StandardService",
        "id": "`+id+`",
        "name": "`+name+`",
        "description": "`+description+`",
        "trusted": true,
        "tags": [
          "string"
        ],
        "serviceUrls": [
          {
            "url": "string",
            "sortOrder": 0,
            "authenticationMethodId": {
              "reference": "string"
            }
          }
        ],
        "meta": {
          "lastUpdated": "2021-11-09T22:18:19.421Z",
          "versionId": "string"
        },
        "organizationGuid": {
          "system": "string",
          "value": "`+organizationID+`"
        }
      }
    }
  ]
}`)
		}
	})
	muxMDM.HandleFunc("/connect/mdm/StandardService/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdService)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdService)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var c mdm.StandardService
	c.Name = name
	c.Description = description
	c.OrganizationGuid = mdm.Identifier{
		Value: organizationID,
	}

	created, resp, err := mdmClient.StandardServices.CreateStandardService(c)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, name, created.Name)

	created, resp, err = mdmClient.StandardServices.GetStandardServiceByID(created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdService) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, id, created.ID)

	ok, resp, err := mdmClient.StandardServices.DeleteStandardService(*created)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdService)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
