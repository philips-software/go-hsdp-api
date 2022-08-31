package dicom_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
)

func TestObjectStoreCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f5a1e608-6787-4af1-a94a-4cbda7677a9c"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/objectStores", func(w http.ResponseWriter, r *http.Request) {
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
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/objectStores/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		switch r.Method {
		case "GET":
			store := dicom.ObjectStore{
				ID:          storeID,
				Description: "Foo",
				StaticAccess: &dicom.StaticAccess{
					Endpoint:   "https://foo",
					BucketName: "bucket-id",
					SecretKey:  "swanson",
					AccessKey:  "ron",
				},
			}
			resp, err := json.Marshal(&store)
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

	created, resp, err := dicomClient.Config.CreateObjectStore(dicom.ObjectStore{
		Description: "Test Store",
		AccessType:  "static",
		StaticAccess: &dicom.StaticAccess{
			Endpoint:   "https://foo",
			BucketName: "bucket-id",
			SecretKey:  "swanson",
			AccessKey:  "ron",
		},
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
	store, resp, err := dicomClient.Config.GetObjectStore(storeID, nil)
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
	ok, resp, err := dicomClient.Config.DeleteObjectStore(dicom.ObjectStore{ID: storeID}, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)
}
