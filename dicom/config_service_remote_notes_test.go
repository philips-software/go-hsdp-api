package dicom_test

import (
	"encoding/json"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestRemoteNodesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	nodeID := "f65b7642-442e-4597-a64f-260f9251ca1d"
	orgID := "614b0053-7a57-44d8-ba8a-809b9362a9a6"

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/remoteNodes", func(w http.ResponseWriter, r *http.Request) {
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
			var received dicom.RemoteNode
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			received.ID = nodeID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			repos := []dicom.RemoteNode{
				{
					ID:    nodeID,
					Title: "A title here",
					NetworkConnection: dicom.NetworkConnection{
						Port:     31337,
						HostName: "foo.com",
					},
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
	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/remoteNodes/"+nodeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		switch r.Method {
		case "GET":
			repo := dicom.Repository{
				ID:                  nodeID,
				OrganizationID:      orgID,
				ActiveObjectStoreID: nodeID,
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

	created, resp, err := dicomClient.Config.CreateRemoteNode(dicom.RemoteNode{
		Title: "Some Title here",
		NetworkConnection: dicom.NetworkConnection{
			Port:     31337,
			HostName: "foo.com",
		},
		AETitle: "AE Title here",
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
	assert.Equal(t, created.ID, nodeID)
	repos, resp, err := dicomClient.Config.GetRemoteNodes(&dicom.GetOptions{OrganizationID: &orgID}, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, repos) {
		return
	}
	assert.Equal(t, (*repos)[0].ID, nodeID)
	ok, resp, err := dicomClient.Config.DeleteRemoteNode(dicom.RemoteNode{ID: nodeID}, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)
}
