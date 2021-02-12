package dicom_test

import (
	"encoding/json"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestCDRServiceAccountGetSet(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceAccountID := "f5a1e608-6787-4af1-a94a-4cbda7677a9c"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/cdrServiceAccount", func(w http.ResponseWriter, r *http.Request) {
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
			var received dicom.CDRServiceAccount
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			received.ID = serviceAccountID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			account := dicom.CDRServiceAccount{
				ID: serviceAccountID,
			}
			resp, err := json.Marshal(&account)
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

	created, resp, err := dicomClient.Config.SetCDRServiceAccount(dicom.CDRServiceAccount{
		ServiceID:  "my@service.id.host",
		PrivateKey: "APrivateKeyHere",
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, created.ID, serviceAccountID)

	account, resp, err := dicomClient.Config.GetCDRServiceAccount(nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, account) {
		return
	}
	assert.Equal(t, serviceAccountID, account.ID)
}
