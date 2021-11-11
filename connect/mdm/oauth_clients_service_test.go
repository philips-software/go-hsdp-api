package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"

	"github.com/stretchr/testify/assert"
)

func TestClientCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	clientName := "TestClient"
	clientDescription := "Group description"
	applicationID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	clientID := "TestClient"
	globalReferenceID := "c3fe79e6-13c2-48c1-adfa-826a01d4b31c"
	createdClient := `{
  "resourceType": "OAuthClient",
  "id": "` + id + `",
  "name": "` + clientName + `",
  "description": "string",
  "applicationId": {
    "reference": "string"
  },
  "globalReferenceId": "` + globalReferenceID + `",
  "redirectionURIs": [
    "string"
  ],
  "responseTypes": [
    "string"
  ],
  "userClient": true,
  "bootstrapClientGuid": {
    "system": "string",
    "value": "string"
  },
  "bootstrapClientId": "string",
  "bootstrapClientSecret": "string",
  "bootstrapClientRevoked": true,
  "clientGuid": {
    "system": "string",
    "value": "string"
  },
  "clientId": "` + clientID + `",
  "clientSecret": "string",
  "clientRevoked": true,
  "meta": {
    "lastUpdated": "2021-11-09T20:26:33.442Z",
    "versionId": "string"
  }
}`
	muxMDM.HandleFunc("/connect/mdm/OAuthClient", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, createdClient)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"total": 1,
				"entry": [
					{
						"clientId": "`+clientName+`",
						"type": "Public",
						"name": "`+clientName+`",
						"realms": [
							"/"
						],
						"description": "Device client1",
						"redirectionURIs": [
							"https://something/OAuth2/something"
						],
						"applicationId": "`+applicationID+`",
						"responseTypes": [
							"code id_token",
							"id_token"
						],
						"globalReferenceId": "`+globalReferenceID+`",
						"defaultScopes": [
							"cn"
						],
						"scopes": [
							"mail",
							"sn"
						],
						"disabled": false,
						"id": "`+clientID+`",
						"meta": {
							"versionId": "0",
							"lastModified": "2015-07-29T15:42:03.123Z"
						}
					}
				]
			}`)
		}
	})
	muxMDM.HandleFunc("/connect/mdm/OAuthClient/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdClient)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
    "id": "`+clientID+`",
    "meta": {
      "versionId": "4",
      "lastModified": "2015-07-29T15:42:03.123Z"
    },
    "clientId": "test1",
    "name":"TestClient1",
    "type":"Public",
    "description": "Device client1",
    "redirectionURIs": [
        "https://example.com/please/send/code_here" ],
    "responseTypes" :["code id_token","id_token"],
    "defaultScopes":["cn"],
    "scopes":["sn","cn"],
    "applicationId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "disabled":false,
    "globalReferenceId": "string",
    "consentImplied": false
}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	muxMDM.HandleFunc("/connect/mdm/OAuthClient/"+clientID+"/$scopes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var c mdm.OAuthClient
	c.Name = clientName
	c.Description = clientDescription
	c.ClientID = clientID
	c.GlobalReferenceID = globalReferenceID
	c.BootstrapClientGuid = mdm.Identifier{
		Value: "foo",
	}
	c.ClientGuid = &mdm.Identifier{
		Value: "bar",
	}

	created, resp, err := mdmClient.OAuthClients.CreateOAuthClient(c)
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
	assert.Equal(t, clientName, created.Name)

	created, resp, err = mdmClient.OAuthClients.GetOAuthClientByID(created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdClient) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, id, created.ID)

	ok, resp, err := mdmClient.OAuthClients.DeleteOAuthClient(*created)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdClient)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
