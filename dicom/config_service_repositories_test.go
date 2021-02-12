package dicom_test

import (
	"encoding/json"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestRepositoriesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f65b7642-442e-4597-a64f-260f9251ca1d"
	orgID := "614b0053-7a57-44d8-ba8a-809b9362a9a6"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/dicomRepositories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		switch r.Method {
		case "POST":
			if !assert.Equal(t, "application/json", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, dicom.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			var receivedStore dicom.ObjectStore
			err := json.NewDecoder(r.Body).Decode(&receivedStore)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			receivedStore.ID = storeID
			resp, err := json.Marshal(&receivedStore)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			repos := []dicom.Repository{
				{
					ID:                  storeID,
					OrganizationID:      orgID,
					ActiveObjectStoreID: storeID,
				},
			}
			resp, err := json.Marshal(&repos)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, string(resp))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/dicomRepositories/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		switch r.Method {
		case "GET":
			repo := dicom.Repository{
				ID:                  storeID,
				OrganizationID:      orgID,
				ActiveObjectStoreID: storeID,
			}
			resp, err := json.Marshal(&repo)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, string(resp))
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := dicomClient.Config.CreateRepository(dicom.Repository{
		OrganizationID:      orgID,
		ActiveObjectStoreID: storeID,
	}, &dicom.QueryOptions{OrganizationID: &orgID})
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
	assert.Equal(t, created.ID, storeID)
	repos, resp, err := dicomClient.Config.GetRepositories(&dicom.QueryOptions{OrganizationID: &orgID}, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, repos) {
		return
	}
	assert.Equal(t, (*repos)[0].ID, storeID)
	ok, resp, err := dicomClient.Config.DeleteRepository(dicom.Repository{ID: storeID}, &dicom.QueryOptions{OrganizationID: &orgID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)
}
