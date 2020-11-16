package iam

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserSelfRegistration(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	newUserUUID := "867128a6-0e02-431c-ba1e-9e764436dae4"
	loginID := "loafoe"
	email := "foo@bar.com"
	firstName := "La"
	lastName := "Foe"

	muxIDM.HandleFunc("/authorize/identity/User", func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("No Authorization header expected, Got: %s", auth)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.Method {
		case "POST":
			var person Person
			body, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(body, &person)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if person.IsAgeValidated != "true" {
				t.Errorf("ageValidated should be true")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.Header().Set("Location", "/authorize/identity/User/"+newUserUUID)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
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
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "total": 1,
  "entry": [
    {
      "preferredLanguage": "en-US",
      "loginId": "`+loginID+`",
      "emailAddress": "`+email+`",
      "id": "`+newUserUUID+`",
      "managingOrganization": "c29cdb88-7cda-4fc1-af8b-ee5947659958",
      "name": {
        "given": "`+firstName+`",
        "family": "`+lastName+`"
      },
      "memberships": [
        {
          "organizationId": "c29cdb88-7cda-4fc1-af8b-ee5947659958",
          "organizationName": "Pawnee",
          "roles": [
            "ANALYZE",
            "S3CREDSADMINROLE",
            "ADMIN"
          ],
          "groups": [
            "S3CredsAdminGroup",
            "AdminGroup"
          ]
        },
        {
          "organizationId": "d4be75cc-e81b-4d7d-b034-baf6f3f10792",
          "organizationName": "Eagleton",
          "roles": [
            "LOGUSER"
          ],
          "groups": [
            "LogUserGroup"
          ]
        }
      ],
      "passwordStatus": {
        "passwordExpiresOn": "2022-02-04T10:07:55Z",
        "passwordChangedOn": "2020-02-15T10:07:55Z"
      },
      "accountStatus": {
        "mfaStatus": "NOTREQUIRED",
        "lastLoginTime": "2020-05-09T12:27:41Z",
        "emailVerified": true,
        "numberOfInvalidAttempt": 0,
        "disabled": false
      },
      "consentedApps": [
        "default default default"
      ]
    }
  ]
}`)
		}

	})
	person := Person{
		ResourceType: "Person",
		LoginID:      loginID,
		Name: Name{
			Family: lastName,
			Given:  firstName,
		},
		Telecom: []TelecomEntry{
			{
				System: "mobile",
				Value:  "3112345678",
			},
			{
				System: "email",
				Value:  "john@doe.com",
			},
		},
		IsAgeValidated: "true",
	}
	user, resp, err := client.Users.CreateUser(person)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	if !assert.NotNil(t, user) {
		return
	}
	assert.Equal(t, newUserUUID, user.ID)
	assert.Equal(t, loginID, user.LoginID)
}

func TestDeleteUser(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	userUUID := "2eec7b01-1417-4546-9c5e-088dea0a9e8b"

	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("No Authorization header expected, Got: %s", auth)
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusNoContent)
	})
	person := Person{
		ID:           userUUID,
		ResourceType: "Person",
		LoginID:      "loafoe",
		Name: Name{
			Family: "Foe",
			Given:  "La",
		},
		Telecom: []TelecomEntry{
			{
				System: "mobile",
				Value:  "3112345678",
			},
			{
				System: "email",
				Value:  "john@doe.com",
			},
		},
		IsAgeValidated: "true",
	}
	ok, resp, err := client.Users.DeleteUser(person)
	if assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGetUsers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	groupID := "1eec7b01-1417-4546-9c5e-088dea0a9e8b"

	muxIDM.HandleFunc("/security/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
			return
		}
		qp := r.URL.Query()
		if ps := qp.Get("pageSize"); ps != "" && ps != "5" {
			t.Errorf("Expected pageSize to be 5, Got: %s", ps)
			return
		}
		if pn := qp.Get("pageNumber"); pn != "" && pn != "1" {
			t.Errorf("Expected pageNumber to be 1, Got: %s", pn)
			return
		}
		_, _ = io.WriteString(w, `{
			"exchange": {
				"users": [
					{
						"userUUID": "7dbfe5fc-1320-4bc6-92a7-2be5d7f07cac"
					},
					{
						"userUUID": "5620b687-7f67-4222-b7c2-91ff312b3066"
					},
					{
						"userUUID": "41c79d7f-c078-4288-8f6d-459292858f00"
					},
					{
						"userUUID": "beba9f50-22ad-4637-ac00-404d8eae4f9d"
					},
					{
						"userUUID": "461ce8d0-7aab-4982-9e1f-cfafb69e51f0"
					}
				],
				"nextPageExists": true
			},
			"responseCode": "200",
			"responseMessage": "Success"
		}`)
	})

	pageNumber := "1"
	pageSize := "5"
	list, resp, err := client.Users.GetUsers(&GetUserOptions{
		GroupID:    &groupID,
		PageNumber: &pageNumber,
		PageSize:   &pageSize,
	})
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, list) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 5, len(list.UserUUIDs))
	assert.True(t, list.HasNextPage)

	foundUser, resp, err := client.Users.LegacyGetUserIDByLoginID("foo")
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "7dbfe5fc-1320-4bc6-92a7-2be5d7f07cac", foundUser)
}

func userIDByLoginIDHandler(t *testing.T, loginID, email, userUUID string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		userId := r.URL.Query().Get("userId")

		if !(userId == loginID || userId == userUUID) {
			_, _ = io.WriteString(w, `{
				"total": 0,
				"entry": []
			}`)
			return
		}
		_, _ = io.WriteString(w, `{
  "total": 1,
  "entry": [
    {
      "preferredLanguage": "en-US",
      "loginId": "`+loginID+`",
      "emailAddress": "`+email+`",
      "id": "`+userUUID+`",
      "managingOrganization": "c29cdb88-7cda-4fc1-af8b-ee5947659958",
      "name": {
        "given": "Ron",
        "family": "Swanson"
      },
      "memberships": [
        {
          "organizationId": "c29cdb88-7cda-4fc1-af8b-ee5947659958",
          "organizationName": "Pawnee",
          "roles": [
            "ANALYZE",
            "S3CREDSADMINROLE",
            "ADMIN"
          ],
          "groups": [
            "S3CredsAdminGroup",
            "AdminGroup"
          ]
        },
        {
          "organizationId": "d4be75cc-e81b-4d7d-b034-baf6f3f10792",
          "organizationName": "Eagleton",
          "roles": [
            "LOGUSER"
          ],
          "groups": [
            "LogUserGroup"
          ]
        }
      ],
      "passwordStatus": {
        "passwordExpiresOn": "2022-02-04T10:07:55Z",
        "passwordChangedOn": "2020-02-15T10:07:55Z"
      },
      "accountStatus": {
        "mfaStatus": "NOTREQUIRED",
        "lastLoginTime": "2020-05-09T12:27:41Z",
        "emailVerified": true,
        "numberOfInvalidAttempt": 0,
        "disabled": false
      },
      "consentedApps": [
        "default default default"
      ]
    }
  ]
}`)
	}
}

func TestGetUserIDByLoginID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	loginID := "ron"
	email := "foo@bar.com"
	muxIDM.HandleFunc("/authorize/identity/User", userIDByLoginIDHandler(t, loginID, email, userUUID))

	uuid, resp, err := client.Users.GetUserIDByLoginID(loginID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, userUUID, uuid)
}

func TestGetUserByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userUUID := "44d20214-7879-4e35-923d-f9d4e01c9746"
	loginID := "ron"
	email := "foo@bar.com"

	muxIDM.HandleFunc("/authorize/identity/User", userIDByLoginIDHandler(t, loginID, email, userUUID))

	foundUser, resp, err := client.Users.GetUserByID(userUUID)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, foundUser) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, email, foundUser.EmailAddress)
	assert.Equal(t, "Swanson", foundUser.Name.Family)
}

func TestUserActions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxIDM.HandleFunc("/authorize/identity/User/$resend-activation",
		actionRequestHandler(t, "resendOTP", "Password reset send", http.StatusOK))
	muxIDM.HandleFunc("/authorize/identity/User/$set-password",
		actionRequestHandler(t, "setPassword", "TODO: fix", http.StatusOK))
	muxIDM.HandleFunc("/authorize/identity/User/$change-password",
		actionRequestHandler(t, "changePassword", "TODO: fix", http.StatusOK))
	muxIDM.HandleFunc("/authorize/identity/User/$recover-password",
		actionRequestHandler(t, "recoverPassword", "TODO: fix", http.StatusOK))
	userUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	loginID := "ron"
	email := "foo@bar.com"
	muxIDM.HandleFunc("/authorize/identity/User", userIDByLoginIDHandler(t, loginID, email, userUUID))
	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID+"/$mfa",
		actionRequestHandler(t, "setMFA", "TODO: fix", http.StatusAccepted))
	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID+"/$unlock",
		actionRequestHandler(t, "unlock", "", http.StatusNoContent))
	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID+"/$change-loginid",
		actionRequestHandler(t, "unlock", "", http.StatusNoContent))

	ok, resp, err := client.Users.ResendActivation("foo@bar.com")
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Users.RecoverPassword("foo@bar.co")
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Users.ChangePassword("foo@bar.co", "0ld", "N3w")
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Users.SetPassword("foo@bar.com", "1234", "newp@ss", "context")
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	uuid, resp, err := client.Users.GetUserIDByLoginID(loginID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, userUUID, uuid)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Users.SetMFAByLoginID(loginID, true)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	ok, resp, err = client.Users.Unlock(userUUID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	ok, resp, err = client.Users.ChangeLoginID(Person{
		ID:      userUUID,
		LoginID: "ronswanon1",
	}, "ronswanson2")
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func actionRequestHandler(t *testing.T, paramName, informationalMessage string, statusCode int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("No Authorization header expected, Got: %s", auth)
		}
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(statusCode)
		if informationalMessage != "" {
			_, _ = io.WriteString(w, `{
			"resourceType": "OperationOutcome",
			"issue": [
				{
					"severity": "information",
					"code": "informational",
					"details": {
						"text": "`+informationalMessage+`"
					}
				}
			]
		  }`)
		}
	}
}
