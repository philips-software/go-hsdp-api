package inference_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/stretchr/testify/assert"
)

func TestComputeTargetCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	targetID := "3a5fa898-0e95-4453-a399-97a65a1bbaf9"

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/ComputeTarget/"+targetID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			data := ai.ComputeTarget{
				ResourceType: "ComputeTarget",
				ID:           targetID,
				Name:         "TestTarget",
				Description:  "TestDescription",
				InstanceType: "ml.p3.16xlarge",
				Storage:      20,
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

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/ComputeTarget", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var received ai.ComputeTarget
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			received.Created = "2021-08-31 16:42:00"
			received.CreatedBy = "test"
			received.ID = targetID
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

	target := ai.ComputeTarget{
		Name:         "TestTarget",
		Description:  "TesTDescription",
		InstanceType: "ml.p3.16xlarge",
		Storage:      20,
	}

	created, _, err := inferenceClient.ComputeTarget.CreateComputeTarget(target)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, targetID, created.ID)

	retrieved, _, err := inferenceClient.ComputeTarget.GetComputeTargetByID(targetID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, retrieved) {
		return
	}
	assert.Equal(t, targetID, retrieved.ID)

	resp, err := inferenceClient.ComputeTarget.DeleteComputeTarget(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
