package cdl_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/stretchr/testify/assert"
)

func TestResearchStudiesCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	studyOwnerID := "f65b7642-442e-4597-a64f-260f9251ca1d"
	studyID := "614b0053-7a57-44d8-ba8a-809b9362a9a6"
	studyTitle := "My test study"
	studyDescription := "My Test study description"

	someStudy := cdl.Study{
		Title:       studyTitle,
		Description: studyDescription,
		StudyOwner:  studyOwnerID,
	}
	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/Study/"+studyID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET", "PUT":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
  "id": "`+studyID+`",
  "title": "`+studyTitle+`",
  "description": "`+studyDescription+`",
  "organization": "`+cdlTenantID+`",
  "studyOwner": "`+studyOwnerID+`",
  "period": {
    "end": "2022-01-01T00:00:00.000Z"
  },
  "dataProtectedFromDeletion": false
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/Study", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			if !assert.Equal(t, "application/json", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, notification.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			var received cdl.Study
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			received.ID = studyID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
  "resourceType": "Bundle",
  "id": "38c0fc92-a083-412f-a04d-b0b5035e7d88",
  "type": "searchset",
  "link": [
    {
      "relation": "next",
      "url": "https://cicd-datalake.cloud.pcftest.com/store/cdl/`+cdlTenantID+`/Study?page=1"
    },
    {
      "relation": "self",
      "url": "https://cicd-datalake.cloud.pcftest.com/store/cdl/`+cdlTenantID+`/Study"
    }
  ],
  "entry": [
    {
      "fullUrl": "https://cicd-datalake.cloud.pcftest.com/store/cdl/`+cdlTenantID+`/Study/`+studyID+`",
      "resource": {
        "resourceType": "ResearchStudy",
        "id": "`+studyID+`",
        "title": "`+studyTitle+`",
        "description": "`+studyDescription+`",
        "organization": "`+cdlTenantID+`",
        "studyOwner": "`+studyOwnerID+`",
        "dataProtectedFromDeletion": false
      }
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := cdlClient.Study.CreateStudy(someStudy)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, created.ID, studyID)

	item, resp, err := cdlClient.Study.GetStudyByID(studyID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, studyID, item.ID)

	item, resp, err = cdlClient.Study.UpdateStudy(*item)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, studyID, item.ID)
}
