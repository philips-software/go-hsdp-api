package inference

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type JobService struct {
	client *Client

	validate *validator.Validate
}

type InputEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type OutputEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ReferenceComputeModel struct {
	Reference  string `json:"reference"`
	Identifier string `json:"identifier,omitempty"`
}

type ReferenceComputeTarget struct {
	Reference  string `json:"reference"`
	Identifier string `json:"identifier,omitempty"`
}

type Job struct {
	ID                      string                 `json:"id,omitempty"`
	ResourceType            string                 `json:"resourceType"`
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	Type                    string                 `json:"type"`
	Model                   ReferenceComputeModel  `json:"model"`
	ComputeTarget           ReferenceComputeTarget `json:"computeTarget"`
	Input                   []InputEntry           `json:"input"`
	Output                  []OutputEntry          `json:"output"`
	EnvVars                 []EnvironmentVariable  `json:"envVars"`
	CommandArgs             []string               `json:"commandArgs"`
	Status                  string                 `json:"status,omitempty"`
	StatusMessage           string                 `json:"statusMessage,omitempty"`
	Labels                  []string               `json:"labels,omitempty"`
	CreatedBy               string                 `json:"createdBy,omitempty"`
	Created                 string                 `json:"created,omitempty"`
	Completed               string                 `json:"completed,omitempty"`
	Duration                int                    `json:"duration,omitempty"`
	Timeout                 int                    `json:"timeOut,omitempty"`
	AdditionalConfiguration string                 `json:"additionalConfiguration,omitempty"`
}

func (s *JobService) path(components ...string) string {
	return path.Join(components...)
}

func (s *JobService) CreateJob(job Job) (*Job, *Response, error) {
	if err := s.validate.Struct(job); err != nil {
		return nil, nil, err
	}
	req, err := s.client.newInferenceRequest("POST", s.path("InferenceJob"), job, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var createdJob Job
	resp, err := s.client.do(req, &createdJob)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateJob: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdJob, resp, nil
}

func (s *JobService) DeleteJob(job Job) (*Response, error) {
	req, err := s.client.newInferenceRequest("DELETE", s.path("InferenceJob", job.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	resp, err := s.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteJob: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *JobService) GetJobByID(id string) (*Job, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("InferenceJob", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var foundJob Job
	resp, err := s.client.do(req, &foundJob)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetJobByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundJob, resp, nil
}

func (s *JobService) GetJobs(opt *GetOptions, options ...OptionFunc) ([]Job, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("InferenceJob"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var bundleResponse struct {
		ResourceType string                 `json:"resourceType,omitempty"`
		Type         string                 `json:"type,omitempty"`
		Total        int                    `json:"total,omitempty"`
		Entry        []internal.BundleEntry `json:"entry"`
	}
	resp, err := s.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	var jobs []Job
	for _, e := range bundleResponse.Entry {
		var job Job
		if err := json.Unmarshal(e.Resource, &job); err == nil {
			jobs = append(jobs, job)
		}
	}
	return jobs, resp, err
}

func (s *JobService) TerminateJob(job Job) (*Response, error) {
	req, err := s.client.newInferenceRequest("POST", s.path("InferenceJob", job.ID, "$terminate"), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	resp, err := s.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("TerminateJob: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}
