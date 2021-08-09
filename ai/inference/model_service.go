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

type ModelService struct {
	client *Client

	validate *validator.Validate
}

type ReferenceComputeEnvironment struct {
	Reference  string `json:"reference"`
	Identifier string `json:"identifier,omitempty"`
}

type Sourcecode struct {
	URL      string `json:"url" validate:"required"`
	Branch   string `json:"branch"`
	CommitID string `json:"commitID"`
	SSHKey   string `json:"sshKey"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Model struct {
	ID                      string                      `json:"id,omitempty"`
	ResourceType            string                      `json:"resourceType"`
	Name                    string                      `json:"name" validate:"required"`
	Version                 string                      `json:"version" validate:"required"`
	Description             string                      `json:"description"`
	ComputeEnvironment      ReferenceComputeEnvironment `json:"computeEnvironment" validate:"required"`
	ArtifactPath            string                      `json:"artifactPath"`
	Sourcecode              Sourcecode                  `json:"sourceCode"`
	EntryCommands           []string                    `json:"entryCommands" validate:"required"`
	EnvVars                 []EnvironmentVariable       `json:"envVars,omitempty"`
	Labels                  []string                    `json:"labels"`
	Type                    string                      `json:"type"`
	AdditionalConfiguration string                      `json:"additionalConfiguration,omitempty"`
	Created                 string                      `json:"created,omitempty"`
	CreatedBy               string                      `json:"createdBy,omitempty"`
}

func (s *ModelService) path(components ...string) string {
	return path.Join(components...)
}

func (s *ModelService) CreateModel(model Model) (*Model, *Response, error) {
	if err := s.validate.Struct(model); err != nil {
		return nil, nil, err
	}
	req, err := s.client.newInferenceRequest("POST", s.path("Model"), model, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var createdModel Model
	resp, err := s.client.do(req, &createdModel)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateModel: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdModel, resp, nil
}

func (s *ModelService) DeleteModel(model Model) (*Response, error) {
	req, err := s.client.newInferenceRequest("DELETE", s.path("Model", model.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	resp, err := s.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteModel: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *ModelService) GetModelByID(id string) (*Model, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("Model", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var foundModel Model
	resp, err := s.client.do(req, &foundModel)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetModelByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundModel, resp, nil
}

func (s *ModelService) GetModels(opt *GetOptions, options ...OptionFunc) ([]Model, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("Model"), opt, options...)
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
	var models []Model
	for _, e := range bundleResponse.Entry {
		var model Model
		if err := json.Unmarshal(e.Resource, &model); err == nil {
			models = append(models, model)
		}
	}
	return models, resp, err
}
