package dicom_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
)

func TestFHIRStoreGetSet(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f5a1e608-6787-4af1-a94a-4cbda7677a9c"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/fhirStore", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
			var receivedStore dicom.FHIRStore
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
			store := dicom.FHIRStore{
				ID: storeID,
			}
			resp, err := json.Marshal(&store)
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

	created, resp, err := dicomClient.Config.SetFHIRStore(dicom.FHIRStore{
		MPIEndpoint: "https://foo.bar/mpi",
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
	assert.Equal(t, created.ID, storeID)

	store, resp, err := dicomClient.Config.GetFHIRStore(nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, store) {
		return
	}
	assert.Equal(t, store.ID, storeID)
}
