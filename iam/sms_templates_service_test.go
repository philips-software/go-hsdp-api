package iam

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTemplate = `{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSTemplate"
  ],
  "id": "63f9d48c-0502-49aa-bc20-b5ca483e275f",
  "organization": {
    "value": "c57b2625-eda3-4b27-a8e6-86f0a0e76afc"
  },
  "type": "PHONE_VERIFICATION",
  "message": "SGVsbG8gd29ybGQ=",
  "locale": "en-US",
  "createdBy": {
    "value": "9b6f1d8a-0967-42d8-9622-d30c877c9da4"
  },
  "modifiedBy": {
    "value": "9b6f1d8a-0967-42d8-9622-d30c877c9da4"
  },
  "meta": {
    "resourceType": "SMSTemplate",
    "created": "2021-10-13T07:46:08.261Z",
    "lastModified": "2021-10-13T07:46:08.261Z",
    "location": "https://foo.bar.com/authorize/scim/v2/Configurations/SMSTemplate/63f9d48c-0502-49aa-bc20-b5ca483e275f",
    "version": "W/\"-151862790\""
  }
}`
)

func TestCreateSMSTemplate(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "c57b2625-eda3-4b27-a8e6-86f0a0e76afc"

	muxIDM.HandleFunc("/authorize/scim/v2/Configurations/SMSTemplate", func(w http.ResponseWriter, r *http.Request) {
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
		var newTemplate SMSTemplate
		err = json.Unmarshal(body, &newTemplate)
		if err != nil {
			t.Errorf("Expected orgnization in body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newTemplate.Organization.Value != orgID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, testTemplate)
	})
	var newTemplate = SMSTemplate{
		Organization: OrganizationValue{
			Value: "c57b2625-eda3-4b27-a8e6-86f0a0e76afc",
		},
		Type:    TypePhoneVerification,
		Locale:  "en-US",
		Message: base64.StdEncoding.EncodeToString([]byte("Hello world")),
	}

	createdTemplate, resp, err := client.SMSTemplates.CreateSMSTemplate(newTemplate)
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, createdTemplate) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, TypePhoneVerification, createdTemplate.Type)
	assert.Equal(t, orgID, createdTemplate.Organization.Value)
}
