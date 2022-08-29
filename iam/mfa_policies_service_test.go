package iam

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMFAPolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	policyID := "400f1adb-bba6-4f52-8d04-f78ecd3833da"
	userID := "ad8a7c6a-231e-452c-8e89-9863c1005982"
	orgID := "b23e7a82-f3b4-40b9-aaef-8111cb788ef9"
	muxIDM.HandleFunc("/authorize/scim/v2/MFAPolicies", func(w http.ResponseWriter, r *http.Request) {
		if ok := assert.Equal(t, "POST", r.Method); !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var newPolicy MFAPolicy
		if err := json.Unmarshal(body, &newPolicy); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newID := policyID
		newType := "User"
		newValue := userID
		if newPolicy.Resource.Type == "Organization" {
			newType = "Organization"
			newValue = orgID
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"schemas": [
			  "urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:MFAPolicy"
			],
			"id": "`+newID+`",
			"name": "TestPolicy",
			"resource": {
			  "type": "`+newType+`",
			  "value": "`+newValue+`"
			},
			"types": [
			  "SOFT_OTP"
			],
			"active": true,
			"createdBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"modifiedBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"meta": {
			  "resourceType": "MFAPolicy",
			  "created": "2019-12-09T07:38:05.183Z",
			  "lastModified": "2019-12-09T07:38:05.183Z",
			  "location": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/scim/v2/MFAPolicies/`+newID+`",
			  "version": "W/\"-955544145\""
			}
		  }`)
	})
	var policy MFAPolicy
	policy.SetResourceUser(userID)
	policy.SetType("SOFT_OTP")
	policy.SetActive(true)

	newPolicy, resp, err := client.MFAPolicies.CreateMFAPolicy(policy)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, newPolicy)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, newPolicy.ID, policyID)
	assert.Equal(t, "User", newPolicy.Resource.Type)
	assert.True(t, *newPolicy.Active)

	var orgPolicy MFAPolicy
	orgPolicy.SetResourceOrganization(orgID)
	orgPolicy.SetType("SOFT_OTP")
	orgPolicy.SetActive(true)
	newPolicy, resp, err = client.MFAPolicies.CreateMFAPolicy(orgPolicy)
	if err != nil {
		t.Fatal(err)
	}
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if ok := assert.NotNil(t, newPolicy); ok {
		assert.Equal(t, newPolicy.ID, policyID)
		assert.Equal(t, "Organization", newPolicy.Resource.Type)
		assert.True(t, *newPolicy.Active)
	}
}

func TestGetMFAPolicyByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	policyID := "400f1adb-bba6-4f52-8d04-f78ecd3833da"
	userID := "ad8a7c6a-231e-452c-8e89-9863c1005982"
	muxIDM.HandleFunc("/authorize/scim/v2/MFAPolicies/"+policyID, func(w http.ResponseWriter, r *http.Request) {
		if ok := assert.Equal(t, "GET", r.Method); !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/scim+json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"schemas": [
			  "urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:MFAPolicy"
			],
			"id": "`+policyID+`",
			"name": "TestPolicy",
			"resource": {
			  "type": "User",
			  "value": "`+userID+`"
			},
			"types": [
			  "SOFT_OTP"
			],
			"active": true,
			"createdBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"modifiedBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"meta": {
			  "resourceType": "MFAPolicy",
			  "created": "2019-12-09T07:38:05.183Z",
			  "lastModified": "2019-12-09T07:38:05.183Z",
			  "location": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/scim/v2/MFAPolicies/`+policyID+`",
			  "version": "W/\"-955544145\""
			}
		  }`)
	})

	policy, resp, err := client.MFAPolicies.GetMFAPolicyByID(policyID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	if ok := assert.NotNil(t, policy); ok {
		assert.Equal(t, policy.ID, policyID)
		assert.True(t, *policy.Active)
	}
}

func TestUpdateMFAPolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	policyID := "400f1adb-bba6-4f52-8d04-f78ecd3833da"
	userID := "ad8a7c6a-231e-452c-8e89-9863c1005982"
	description := "New description"

	muxIDM.HandleFunc("/authorize/scim/v2/MFAPolicies/"+policyID, func(w http.ResponseWriter, r *http.Request) {
		if ok := assert.Equal(t, "PUT", r.Method); !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		etag := r.Header.Get("If-Match")
		if etag != "W/\"-955544145\"" {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
		w.Header().Set("Content-Type", "application/scim+json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"schemas": [
			  "urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:MFAPolicy"
			],
			"id": "`+policyID+`",
			"name": "TestPolicy",
			"resource": {
			  "type": "User",
			  "value": "`+userID+`"
			},
			"types": [
			  "SOFT_OTP"
			],
			"active": true,
			"createdBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"modifiedBy": {
			  "value": "3ad37d5c-77fe-483e-a7f2-535cc5cd397e"
			},
			"meta": {
			  "resourceType": "MFAPolicy",
			  "created": "2019-12-09T07:38:05.183Z",
			  "lastModified": "2019-12-09T07:38:05.183Z",
			  "location": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/scim/v2/MFAPolicies/400f1adb-bba6-4f52-8d04-f78ecd3833da",
			  "version": "W/\"-955544145\""
			}
		  }`)
	})
	var policy MFAPolicy
	policy.ID = policyID
	policy.Description = description
	policy.Meta = &MFAPolicyMeta{
		Version: "W/\"-955544145\"",
	}

	updatedPolicy, resp, err := client.MFAPolicies.UpdateMFAPolicy(&policy)
	assert.Nil(t, err)
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.NotNil(t, updatedPolicy)
}

func TestDeleteMFAPolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	policyID := "400f1adb-bba6-4f52-8d04-f78ecd3833da"
	muxIDM.HandleFunc("/authorize/scim/v2/MFAPolicies/"+policyID, func(w http.ResponseWriter, r *http.Request) {
		if ok := assert.Equal(t, "DELETE", r.Method); !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	var policy = MFAPolicy{ID: policyID}
	ok, resp, err := client.MFAPolicies.DeleteMFAPolicy(policy)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, true, ok, "expected MFA policy deletion to succeed")
}
