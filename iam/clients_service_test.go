package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	clientName := "TestClient"
	clientDescription := "Group description"
	applicationID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	clientID := "TestClient"
	password := "SomePassword"
	globalReferenceID := "c3fe79e6-13c2-48c1-adfa-826a01d4b31c"
	muxIDM.HandleFunc("/authorize/identity/Client", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Location", "/authorize/identity/Client/"+clientID)
			w.Header().Set("ETag", "0")
			w.WriteHeader(http.StatusCreated)
		case "GET":
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
	muxIDM.HandleFunc("/authorize/identity/Client/"+clientID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
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
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Client/"+clientID+"/$scopes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var c ApplicationClient
	c.Name = clientName
	c.Description = clientDescription
	c.ApplicationID = applicationID
	c.ClientID = clientID
	c.Password = password
	c.GlobalReferenceID = globalReferenceID

	createdClient, resp, err := client.Clients.CreateClient(c)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdClient)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, clientName, createdClient.Name)

	createdClient, resp, err = client.Clients.GetClientByID(createdClient.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdClient) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, clientID, createdClient.ID)

	createdClient.Password = password
	cl, resp, err := client.Clients.UpdateClient(*createdClient)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, cl) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Public", cl.Type)

	ok, resp, err := client.Clients.DeleteClient(*createdClient)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdClient)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
