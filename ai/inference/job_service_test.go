package inference_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/stretchr/testify/assert"
)

func TestJobCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	jobID := "3a5fa898-0e95-4453-a399-97a65a1bbaf9"

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/InferenceJob/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			data := ai.Job{
				ResourceType: "InferenceJob",
				ID:           jobID,
				Name:         "TestJob",
				Description:  "TestDescription",
				Timeout:      60,
				Model: ai.ReferenceComputeModel{
					Reference: "Model/ca41fd8e-e4b9-4f52-b065-a6c6671af57b",
				},
				ComputeTarget: ai.ReferenceComputeTarget{
					Reference: "ComputeTarget/d0e9cae2-7563-482b-b53e-a8ae703e101d",
				},
				Type: "sagemaker",
				Labels: []string{
					"CNN",
					"Test",
					"INDEX-0",
				},
				CommandArgs:   []string{"yo"},
				Status:        "Completed",
				StatusMessage: "SomethingSomething",
				EnvVars: []ai.EnvironmentVariable{
					{Name: "Foo", Value: "Bar"},
				},
				Created:   "2021-08-31 16:42:00",
				CreatedBy: "test",
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

	muxInference.HandleFunc("/analyze/inference/"+inferenceTenantID+"/InferenceJob", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var received ai.Job
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			received.Created = "2021-08-31 16:42:00"
			received.CreatedBy = "test"
			received.ID = jobID
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

	data := ai.Job{
		ResourceType: "Model",
		Name:         "TestJob",
		Description:  "TestDescription",
		Timeout:      60,
		Model: ai.ReferenceComputeModel{
			Reference: "Model/ca41fd8e-e4b9-4f52-b065-a6c6671af57b",
		},
		ComputeTarget: ai.ReferenceComputeTarget{
			Reference: "ComputeTarget/d0e9cae2-7563-482b-b53e-a8ae703e101d",
		},
		Type: "sagemaker",
		Labels: []string{
			"CNN",
			"Test",
			"INDEX-0",
		},
		CommandArgs:   []string{"yo"},
		Status:        "Completed",
		StatusMessage: "SomethingSomething",
		EnvVars: []ai.EnvironmentVariable{
			{Name: "Foo", Value: "Bar"},
		},
	}

	created, _, err := inferenceClient.Job.CreateJob(data)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, jobID, created.ID)

	retrieved, _, err := inferenceClient.Job.GetJobByID(jobID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, retrieved) {
		return
	}
	assert.Equal(t, jobID, retrieved.ID)

	resp, err := inferenceClient.Job.DeleteJob(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
