package inference

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/internal"
)

type ModelService struct {
	client *ai.Client

	validate *validator.Validate
}

type Model struct {
	ID                      string                         `json:"id,omitempty"`
	ResourceType            string                         `json:"resourceType"`
	Name                    string                         `json:"name" validate:"required"`
	Version                 string                         `json:"version" validate:"required"`
	Description             string                         `json:"description,omitempty"`
	ComputeEnvironment      ai.ReferenceComputeEnvironment `json:"computeEnvironment" validate:"required"`
	ArtifactPath            string                         `json:"artifactPath,omitempty"`
	SourceCode              ai.SourceCode                  `json:"sourceCode"`
	EntryCommands           []string                       `json:"entryCommands" validate:"required"`
	EnvVars                 []ai.EnvironmentVariable       `json:"envVars,omitempty"`
	Labels                  []string                       `json:"labels,omitempty"`
	Type                    string                         `json:"type,omitempty"`
	AdditionalConfiguration string                         `json:"additionalConfiguration,omitempty"`
	Created                 string                         `json:"created,omitempty"`
	CreatedBy               string                         `json:"createdBy,omitempty"`
}

func (s *ModelService) path(components ...string) string {
	return path.Join(components...)
}

func (s *ModelService) CreateModel(model Model) (*Model, *ai.Response, error) {
	if err := s.validate.Struct(model); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewAIRequest("POST", s.path("Model"), model, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var createdModel Model
	resp, err := s.client.Do(req, &createdModel)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateModel: %w", ai.ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdModel, resp, nil
}

func (s *ModelService) DeleteModel(model Model) (*ai.Response, error) {
	req, err := s.client.NewAIRequest("DELETE", s.path("Model", model.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	resp, err := s.client.Do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteModel: %w", ai.ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *ModelService) GetModelByID(id string) (*Model, *ai.Response, error) {
	req, err := s.client.NewAIRequest("GET", s.path("Model", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var foundModel Model
	resp, err := s.client.Do(req, &foundModel)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetModelByID: %w", ai.ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundModel, resp, nil
}

func (s *ModelService) GetModels(opt *ai.GetOptions, options ...ai.OptionFunc) ([]Model, *ai.Response, error) {
	req, err := s.client.NewAIRequest("GET", s.path("Model"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var bundleResponse struct {
		ResourceType string                 `json:"resourceType,omitempty"`
		Type         string                 `json:"type,omitempty"`
		Total        int                    `json:"total,omitempty"`
		Entry        []internal.BundleEntry `json:"entry"`
	}
	resp, err := s.client.Do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			return nil, resp, ai.ErrEmptyResult
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
