package inference_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/stretchr/testify/assert"
)

func TestComputeEnvironmentCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	envID := "3a5fa898-0e95-4453-a399-97a65a1bbaf9"

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/ComputeEnvironment/"+envID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			data := ai.ComputeEnvironment{
				ResourceType: "ComputeEnvironment",
				ID:           envID,
				Name:         "TestEnv",
				Description:  "TestDescription",
				Image:        "foo.bar/sage/thing",
				IsFactory:    false,
				Created:      "2021-08-31 16:42:00",
				CreatedBy:    "test",
			}
			resp, err := json.Marshal(&data)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/ComputeEnvironment", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var received ai.ComputeEnvironment
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			received.Created = "2021-08-31 16:42:00"
			received.CreatedBy = "test"
			received.ID = envID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, string(resp))
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})

	env := ai.ComputeEnvironment{
		Name:        "TestEnv",
		Description: "TesTDescription",
		Image:       "foo.bar/sage/thing",
	}

	created, _, err := inferenceClient.ComputeEnvironment.CreateComputeEnvironment(env)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, envID, created.ID)

	retrieved, _, err := inferenceClient.ComputeEnvironment.GetComputeEnvironmentByID(envID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, retrieved) {
		return
	}
	assert.Equal(t, envID, retrieved.ID)

	resp, err := inferenceClient.ComputeEnvironment.DeleteComputeEnvironment(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
