package iam

import (
	"encoding/base64"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailTemplateCreateDelete(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "2c266886-f918-4223-941d-437cb3cd09e8"
	orgID := "bda40124-54fa-4967-b2fb-23dcc4e0ad1a"

	muxIDM.HandleFunc("/authorize/identity/EmailTemplate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
  "id": "`+id+`",
  "type": "PASSWORD_CHANGED",
  "managingOrganization": "`+orgID+`",
  "format": "HTML",
  "from": "default",
  "subject": "Password changed",
  "message": "WW91ciBwYXNzd29yZCBoYXMgY2hhbmdlZCE=",
  "locale": "en-US",
  "link": "default",
  "meta": {
    "version": "W/\"473187190\"",
    "updatedBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
    "createdBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
    "created": "2021-01-20T06:06:17.332Z",
    "lastModified": "2021-01-20T06:06:17.332Z"
  }
}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "total": 1,
  "entry": [
    {
      "id": "`+id+`",
      "type": "PASSWORD_CHANGED",
      "managingOrganization": "`+orgID+`",
      "format": "HTML",
      "from": "default",
      "subject": "Password changed",
      "locale": "en-US",
      "link": "default",
      "meta": {
        "version": "W/\"473187190\"",
        "updatedBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
        "createdBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
        "created": "2021-01-20T06:06:17.332Z",
        "lastModified": "2021-01-20T06:06:17.332Z"
      }
    }
  ]
}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/EmailTemplate/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "id": "`+id+`",
  "type": "PASSWORD_CHANGED",
  "managingOrganization": "`+orgID+`",
  "format": "HTML",
  "from": "default",
  "subject": "Password changed",
  "locale": "en-US",
  "link": "default",
  "meta": {
    "version": "W/\"473187190\"",
    "updatedBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
    "createdBy": "3bc7880f-1a01-4cc1-babc-95c4e9bb9b5a",
    "created": "2021-01-20T06:06:17.332Z",
    "lastModified": "2021-01-20T06:06:17.332Z"
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	e := EmailTemplate{
		ManagingOrganization: orgID,
		Subject:              "Password changed",
		Type:                 "PASSWORD_CHANGED",
		Format:               "HTML",
		Locale:               "en-US",
		Message:              base64.StdEncoding.EncodeToString([]byte(`Your password has changed!`)),
	}

	template, resp, err := client.EmailTemplates.CreateTemplate(e)
	if ok := assert.Nil(t, err); !ok {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, template) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, id, template.ID)

	foundTemplate, resp, err := client.EmailTemplates.GetTemplateByID(template.ID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, foundTemplate) {
		return
	}
	assert.Equal(t, template.ID, foundTemplate.ID)

	template, resp, err = client.EmailTemplates.GetTemplate(&GetEmailTemplatesOptions{
		OrganizationID: &orgID,
	})
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, template) {
		return
	}

	ok, resp, err := client.EmailTemplates.DeleteTemplate(*template)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.True(t, ok)
}
