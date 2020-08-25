package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordPolicyCrud(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "2c266886-f918-4223-941d-437cb3cd09e8"
	orgID := "bda40124-54fa-4967-b2fb-23dcc4e0ad1a"

	muxIDM.HandleFunc("/authorize/identity/PasswordPolicy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
      "managingOrganization": "`+orgID+`",
      "expiryPeriodInDays": 720,
      "historyCount": 5,
      "complexity": {
        "minLength": 8,
        "maxLength": 16,
        "minNumerics": 1,
        "minUpperCase": 1,
        "minLowerCase": 1,
        "minSpecialChars": 1
      },
      "id": "`+id+`",
      "meta": {
        "version": "W/\"233552990\"",
        "updatedBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "createdBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "created": "2020-04-23T19:40:31.463Z",
        "lastModified": "2020-04-24T20:47:27.473Z"
      }
    }`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"total": 1,
				"entry": [ {
      "managingOrganization": "`+orgID+`",
      "expiryPeriodInDays": 720,
      "historyCount": 5,
      "complexity": {
        "minLength": 8,
        "maxLength": 16,
        "minNumerics": 1,
        "minUpperCase": 1,
        "minLowerCase": 1,
        "minSpecialChars": 1
      },
      "id": "`+id+`",
      "meta": {
        "version": "W/\"233552990\"",
        "updatedBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "createdBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "created": "2020-04-23T19:40:31.463Z",
        "lastModified": "2020-04-24T20:47:27.473Z"
      }
    }
]
	}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/PasswordPolicy/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		case "GET", "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
      "managingOrganization": "`+orgID+`",
      "expiryPeriodInDays": 720,
      "historyCount": 5,
      "complexity": {
        "minLength": 8,
        "maxLength": 16,
        "minNumerics": 1,
        "minUpperCase": 1,
        "minLowerCase": 1,
        "minSpecialChars": 1
      },
      "id": "`+id+`",
      "meta": {
        "version": "W/\"233552990\"",
        "updatedBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "createdBy": "29c58b06-fdaa-461f-afef-53c91a18acbd",
        "created": "2020-04-23T19:40:31.463Z",
        "lastModified": "2020-04-24T20:47:27.473Z"
      }
    }`)

		}
	})

	var p PasswordPolicy
	p.ManagingOrganization = orgID
	p.ExpiryPeriodInDays = 90
	p.HistoryCount = 10

	policy, resp, err := client.PasswordPolicies.CreatePasswordPolicy(p)
	if ok := assert.Nil(t, err); !ok {
		return
	}
	if ok := assert.NotNil(t, resp); !ok {
		return
	}
	if ok := assert.NotNil(t, policy); !ok {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, id, policy.ID)

	foundPolicy, resp, err := client.PasswordPolicies.GetPasswordPolicyByID(policy.ID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, foundPolicy) {
		return
	}
	assert.Equal(t, policy.ID, foundPolicy.ID)

	policies, resp, err := client.PasswordPolicies.GetPasswordPolicies(&GetPasswordPolicyOptions{
		OrganizationID: &orgID,
	})
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, policies) {
		return
	}
	assert.Equal(t, 1, len(*policies))

	updatedPolicy, resp, err := client.PasswordPolicies.UpdatePasswordPolicy(*foundPolicy)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, updatedPolicy) {
		return
	}

	ok, resp, err := client.PasswordPolicies.DeletePasswordPolicy(*foundPolicy)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.True(t, ok)
}
