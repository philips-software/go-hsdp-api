package credentials

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	productKey := "430deb9e-01c8-4a3b-81dd-e2e46569cd5e"
	policyJSON := `{
		"allowed": {
		  "resources": [
			"${managingOrganization}/folder1/*",
			"54ba7674-8722-40b0-95c6-6514083c870e/folder2/*"
		  ],
		  "actions": [
			"PUT"
		  ]
		},
		"conditions": {
		  "managingOrganizations": [
			"d4d84cf0-f5ee-47a1-86e7-db26d679d95e"
		  ],
		  "groups": [
			"PublishGroup"
		  ]
		},
		"id": 1,
		"resourceType": "Policy"
	  }`

	muxCreds.HandleFunc("/core/credentials/Policy", func(w http.ResponseWriter, r *http.Request) {
		if k := r.Header.Get("X-Product-Key"); k != productKey {
			t.Errorf(ErrMissingProductKey.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unexpected EOF from reading request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var policy Policy
		err = json.Unmarshal(body, &policy)
		if err != nil {
			t.Errorf("Expected contract in body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, policyJSON)
	})

	var newPolicy = Policy{}
	err := json.Unmarshal([]byte(policyJSON), &newPolicy)
	assert.Nil(t, err)
	// Reset
	newPolicy.ResourceType = ""
	newPolicy.ID = 0
	// Set Policy
	newPolicy.ProductKey = productKey

	createdPolicy, resp, err := credsClient.Policy.CreatePolicy(newPolicy)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, createdPolicy)
	assert.IsType(t, &Policy{}, createdPolicy)
	assert.Equal(t, 1, createdPolicy.ID)

}

func TestDeletePolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	productKey := "430deb9e-01c8-4a3b-81dd-e2e46569cd5e"
	id := "1"

	muxCreds.HandleFunc("/core/credentials/Policy/"+id, func(w http.ResponseWriter, r *http.Request) {
		if k := r.Header.Get("X-Product-Key"); k != productKey {
			t.Errorf(ErrMissingProductKey.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var newPolicy = Policy{ID: 1, ProductKey: productKey}
	ok, resp, err := credsClient.Policy.DeletePolicy(newPolicy)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, true, ok, "expected policy deletion to succeed")
}

func TestGetPolicy(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	productKey := "430deb9e-01c8-4a3b-81dd-e2e46569cd5e"
	id := "1"
	muxCreds.HandleFunc("/core/credentials/Policy", func(w http.ResponseWriter, r *http.Request) {
		if k := r.Header.Get("X-Product-Key"); k != productKey {
			t.Errorf(ErrMissingProductKey.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("id") != id {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `[]`)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `[
				{
					"allowed": {
					  "resources": [
						"${managingOrganization}/folder1/*",
						"54ba7674-8722-40b0-95c6-6514083c870e/folder2/*"
					  ],
					  "actions": [
						"PUT"
					  ]
					},
					"conditions": {
					  "managingOrganizations": [
						"d4d84cf0-f5ee-47a1-86e7-db26d679d95e"
					  ],
					  "groups": [
						"PublishGroup"
					  ]
					},
					"id": 1,
					"resourceType": "Policy"
				  }
			
			]`)
		}
	})

	intID := 1
	policies, resp, err := credsClient.Policy.GetPolicy(&GetPolicyOptions{ID: &intID, ProductKey: &productKey})
	assert.Nil(t, err)
	assert.NotNil(t, policies)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(policies), "expected one policy")
	assert.Equal(t, "Policy", policies[0].ResourceType)
}
