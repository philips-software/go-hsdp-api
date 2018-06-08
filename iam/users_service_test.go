package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestGetUserIDByLoginID(t *testing.T) {
	teardown := setup()
	defer teardown()

	userUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	muxIDM.HandleFunc("/security/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"exchange": {
				"users": [
					{
						"userUUID": "`+userUUID+`"
					}
				]
			},
			"responseCode": "200",
			"responseMessage": "Success"
		}`)
	})

	uuid, resp, err := client.Users.GetUserIDByLoginID(userUUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if uuid != userUUID {
		t.Errorf("Expected UUID: %s, Got: %s", userUUID, uuid)
	}
}

func TestGetUserByID(t *testing.T) {
	teardown := setup()
	defer teardown()

	userUUID := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxIDM.HandleFunc("/security/users/"+userUUID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"exchange": {
				"loginId": "john.doe@domain.com",
				"profile": {
					"contact": {
						"emailAddress": "john.doe@domain.com"
					},
					"givenName": "John",
					"familyName": "Doe",
					"addresses": [],
					"disabled": false
				}
			},
			"responseCode": "200",
			"responseMessage": "Success"
		}`)
	})

	foundUser, resp, err := client.Users.GetUserByID(userUUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if len(foundUser.Telecom) < 1 {
		t.Errorf("Expected at least one TelecomEntry (email)")
	}
	if foundUser.Telecom[0].Value != "john.doe@domain.com" {
		t.Errorf("Unexpected email: %s", foundUser.Telecom[0].Value)
	}
	if foundUser.Name.Family != "Doe" {
		t.Errorf("Expected family name: %s, got: %s", "Doe", foundUser.Name.Family)
	}
}
