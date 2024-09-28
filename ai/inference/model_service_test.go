package inference_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/ai/inference"
	"github.com/stretchr/testify/assert"
)

func TestModelCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	modelID := "3a5fa898-0e95-4453-a399-97a65a1bbaf9"

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/Model/"+modelID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			data := inference.Model{
				ResourceType: "ComputeTarget",
				ID:           modelID,
				Name:         "TestModel",
				Description:  "TestDescription",
				Version:      "v2",
				ComputeEnvironment: ai.ReferenceComputeEnvironment{
					Reference: "ComputeEnvironment/ca41fd8e-e4b9-4f52-b065-a6c6671af57b",
				},
				SourceCode: ai.SourceCode{
					Branch: "main",
					URL:    "git@github.com/philips-labs/variant.git",
				},
				Type: "sagemaker",
				Labels: []string{
					"CNN",
					"Test",
					"INDEX-0",
				},
				EntryCommands: []string{"variant"},
				Created:       "2021-08-31 16:42:00",
				CreatedBy:     "test",
			}
			resp, err := json.Marshal(&data)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(resp)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/Model", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var received inference.Model
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			received.Created = "2021-08-31 16:42:00"
			received.CreatedBy = "test"
			received.ID = modelID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(resp)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})

	data := inference.Model{
		ResourceType: "Model",
		Name:         "TestModel",
		Description:  "TesTDescription",
		Version:      "v2",
		ComputeEnvironment: ai.ReferenceComputeEnvironment{
			Reference: "ComputeEnvironment/ca41fd8e-e4b9-4f52-b065-a6c6671af57b",
		},
		SourceCode: ai.SourceCode{
			Branch: "main",
			URL:    "git@github.com/philips-labs/variant.git",
		},
		Type: "sagemaker",
		Labels: []string{
			"CNN",
			"Test",
			"INDEX-0",
		},
		EntryCommands: []string{"variant"},
	}

	created, _, err := inferenceClient.Model.CreateModel(data)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, modelID, created.ID)

	retrieved, _, err := inferenceClient.Model.GetModelByID(modelID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, retrieved) {
		return
	}
	assert.Equal(t, modelID, retrieved.ID)

	resp, err := inferenceClient.Model.DeleteModel(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
