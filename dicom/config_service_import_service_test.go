package dicom_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
)

func TestImportServiceGetSet(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceID := "f5a1e608-6787-4af1-a94a-4cbda7677a9c"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/importService", func(w http.ResponseWriter, r *http.Request) {
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
			var receivedStore dicom.ImportService
			err := json.NewDecoder(r.Body).Decode(&receivedStore)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			receivedStore.ID = serviceID
			resp, err := json.Marshal(&receivedStore)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(resp)
		case "GET":
			store := dicom.ImportService{
				ID: serviceID,
			}
			resp, err := json.Marshal(&store)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resp)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := dicomClient.Config.SetImportService(dicom.ImportService{
		AETitle: "What is an AE Title anyway?",
	}, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, created.ID, serviceID)

	store, resp, err := dicomClient.Config.GetImportService(nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, store) {
		return
	}
	assert.Equal(t, store.ID, serviceID)
}
