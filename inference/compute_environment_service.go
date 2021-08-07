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

type ComputeEnvironmentService struct {
	client *Client

	validate *validator.Validate
}

type ComputeEnvironment struct {
	ID           string `json:"id,omitempty"`
	ResourceType string `json:"resourceType" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	Image        string `json:"image" validate:"required"`
	IsFactory    bool   `json:"isFactory,omitempty"`
	Created      string `json:"created,omitempty"`
	CreatedBy    string `json:"createdBy,omitempty"`
}

// GetOptions describes the fields on which you can search for producers
type GetOptions struct {
	Page  *string `url:"_page,omitempty"`
	Count *string `url:"_count,omitempty"`
	Sort  *string `url:"_sort,omitempty"`
}

func (s *ComputeEnvironmentService) path(components ...string) string {
	return path.Join(components...)
}

func (s *ComputeEnvironmentService) CreateComputeEnvironment(env ComputeEnvironment) (*ComputeEnvironment, *Response, error) {
	if err := s.validate.Struct(env); err != nil {
		return nil, nil, err
	}
	req, err := s.client.newInferenceRequest("POST", s.path("ComputeEnvironment"), env, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var createdEnv ComputeEnvironment
	resp, err := s.client.do(req, &createdEnv)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateStudy: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdEnv, resp, nil
}

func (s *ComputeEnvironmentService) DeleteComputeEnvironment(env ComputeEnvironment) (*Response, error) {
	req, err := s.client.newInferenceRequest("DELETE", s.path("ComputeEnvironment", env.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	resp, err := s.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteComputeEnvironment: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *ComputeEnvironmentService) GetComputeEnvironmentByID(id string) (*ComputeEnvironment, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("ComputeEnvironment", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var foundEnv ComputeEnvironment
	resp, err := s.client.do(req, &foundEnv)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetComputeEnvironmentByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundEnv, resp, nil
}

func (s *ComputeEnvironmentService) GetComputeEnvironments(opt *GetOptions, options ...OptionFunc) ([]ComputeEnvironment, *Response, error) {
	req, err := s.client.newInferenceRequest("GET", s.path("ComputeEnvironment"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var bundleResponse struct {
		ResourceType string                 `json:"resourceType,omitempty"`
		Type         string                 `json:"type,omitempty"`
		Entry        []internal.BundleEntry `json:"entry"`
	}
	resp, err := s.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	var envs []ComputeEnvironment
	for _, e := range bundleResponse.Entry {
		var env ComputeEnvironment
		if err := json.Unmarshal(e.Resource, &env); err == nil {
			envs = append(envs, env)
		}
	}
	return envs, resp, err
}
