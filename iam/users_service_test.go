package iam

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Jeffail/gabs"
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
	})
	person := Person{
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
	ok, resp, err := client.Users.CreateUser(person)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
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
		if ps := qp.Get("pageSize"); ps != "5" {
			t.Errorf("Expected pageSize to be 5, Got: %s", ps)
			return
		}
		if pn := qp.Get("pageNumber"); pn != "1" {
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
	assert.Equal(t, 5, len(list.Users))
	assert.True(t, list.HasNextPage)
}

func userIDByLoginIDHandler(t *testing.T, loginID, userUUID string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if r.URL.Query().Get("loginId") != loginID {
			_, _ = io.WriteString(w, `{
				"responseCode": "4010",
				"responseMessage": "User does not exist"
			}`)
			return
		}
		_, _ = io.WriteString(w, `{
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
	}
}

func TestGetUserIDByLoginID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userUUID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	loginID := "foo@bar.com"
	muxIDM.HandleFunc("/security/users", userIDByLoginIDHandler(t, loginID, userUUID))

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

	muxIDM.HandleFunc("/security/users/"+userUUID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
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
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, foundUser) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotEqual(t, 0, len(foundUser.Telecom)) {
		return
	}
	assert.Equal(t, "john.doe@domain.com", foundUser.Telecom[0].Value)
	assert.Equal(t, "Doe", foundUser.Name.Family)
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
	loginID := "foo@bar.com"
	muxIDM.HandleFunc("/security/users", userIDByLoginIDHandler(t, loginID, userUUID))
	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID+"/$mfa",
		actionRequestHandler(t, "setMFA", "TODO: fix", http.StatusAccepted))
	muxIDM.HandleFunc("/authorize/identity/User/"+userUUID+"/$unlock",
		actionRequestHandler(t, "unlock", "TODO: fix", http.StatusNoContent))

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
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, userUUID, uuid)

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
}

func actionRequestHandler(t *testing.T, paramName, informationalMessage string, statusCode int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("No Authorization header expected, Got: %s", auth)
		}
		//body, _ := ioutil.ReadAll(r.Body)
		//j, _ := gabs.ParseJSON(body)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(statusCode)
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
