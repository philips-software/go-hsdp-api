package iam

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testGW = `{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSGateway"
  ],
  "id": "d7d4b8c2-e883-4fe8-9dcc-ccc3b072a1e7",
  "organization": {
    "value": "c57b2625-eda3-4b27-a8e6-86f0a0e76afc"
  },
  "provider": "twilio",
  "properties": {
    "sid": "***REMOVED***",
    "endpoint": "https://api.twilio.com/SubOrg2/sendsms",
    "fromNumber": "+447380336672"
  },
  "credentials": {
    "token": "[sensitive]"
  },
  "activationExpiry": 10,
  "createdBy": {
    "value": "9b6f1d8a-0967-42d8-9622-d30c877c9da4"
  },
  "modifiedBy": {
    "value": "9b6f1d8a-0967-42d8-9622-d30c877c9da4"
  },
  "meta": {
    "resourceType": "SmsGatewayConfiguration",
    "created": "2021-10-12T18:48:11.888Z",
    "lastModified": "2021-10-12T18:48:11.888Z",
    "location": "https://idm-client-test.us-east.philips-healthsuite.com/authorize/scim/v2/Configurations/SMSGateway/d7d4b8c2-e883-4fe8-9dcc-ccc3b072a1e7",
    "version": "W/\"\"-456396008\"\""
  },
  "active": true
}`
)

func TestCreateSMSGateway(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "c57b2625-eda3-4b27-a8e6-86f0a0e76afc"

	muxIDM.HandleFunc("/authorize/scim/v2/Configurations/SMSGateway", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got ‘%s’", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Expected body to be read: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newGW SMSGateway
		err = json.Unmarshal(body, &newGW)
		if err != nil {
			t.Errorf("Expected orgnization in body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newGW.Organization.Value != orgID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, testGW)
	})
	var newGW = SMSGateway{
		Organization: OrganizationValue{
			Value: "c57b2625-eda3-4b27-a8e6-86f0a0e76afc",
		},
		Provider: "twilio",
		Credentials: ProviderCredentials{
			Token: "***REMOVED***",
		},
		Active: true,
		Properties: ProviderProperties{
			SID:        "***REMOVED***",
			Endpoint:   "https://api.twilio.com/SubOrg2/sendsms",
			FromNumber: "+447380336672",
		},
		ActivationExpiry: 10,
	}

	createdGW, resp, err := client.SMSGateways.CreateSMSGateway(newGW)
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, createdGW) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "twilio", createdGW.Provider)
	assert.Equal(t, orgID, createdGW.Organization.Value)
}
