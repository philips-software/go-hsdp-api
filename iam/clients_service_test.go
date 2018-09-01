package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestClientCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	clientName := "TestClient"
	clientDescription := "Group description"
	applicationID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	clientID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
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
			io.WriteString(w, `{
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
							"lastModified": "2018-07-26T18:08:207.010Z"
						}
					}
				]
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Client/"+clientID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
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

	createdClient, resp, err := client.Clients.CreateClient(c)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP OK. Got: %d", resp.StatusCode)
	}
	if createdClient.Name != clientName {
		t.Errorf("Expected Client name: %s, Got: %s", clientName, createdClient.Name)
	}
	createdClient, resp, err = client.Clients.GetClientByID(createdClient.ID)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if createdClient.ID != clientID {
		t.Errorf("Expected to find client with ID: %s, Got: %s", clientID, createdClient.ID)
	}
	ok, resp, err := client.Clients.UpdateScope(*createdClient, []string{"cn", "introspect"}, "cn")
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP no content Got: %d", resp.StatusCode)
	}

	ok, resp, err = client.Clients.DeleteClient(*createdClient)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP no content Got: %d", resp.StatusCode)
	}
	if !ok {
		t.Errorf("Expected client to be deleted")
	}
}
