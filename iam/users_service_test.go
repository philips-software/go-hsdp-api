package iam

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jeffail/gabs"
)

func TestCreateUserSelfRegistration(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxIDM.HandleFunc("/authorize/identity/User", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("No Authorization header expected, Got: %s", auth)
		}
		body, _ := ioutil.ReadAll(r.Body)
		j, _ := gabs.ParseJSON(body)
		ageValidated, ok := j.Path("isAgeValidated").Data().(string)
		if !ok {
			t.Errorf("Missing isAgeValidated field")
		}
		if ageValidated != "true" {
			t.Errorf("ageValidated should be true")
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{
			"resourceType": "OperationOutcome",
			"issue": [
				{
					"severity": "information",
					"code": "informational",
					"details": {
						"text": "User Created Successfully. Verification Email has been sent."
					}
				}
			]
		}`)
	})
	ok, resp, err := client.Users.CreateUser("First", "Last", "andy@foo.com", "31612345678", "")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("Unexpected failure, Got !ok")
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP create, Got: %d", resp.StatusCode)
	}
}

func TestGetUserIDByLoginID(t *testing.T) {
	teardown := setup(t)
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
	teardown := setup(t)
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
