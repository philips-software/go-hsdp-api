package ai

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
	Client *Client

	Validate *validator.Validate
}

type ComputeEnvironment struct {
	ID           string `json:"id,omitempty"`
	ResourceType string `json:"resourceType" Validate:"required"`
	Name         string `json:"name" Validate:"required"`
	Description  string `json:"description"`
	Image        string `json:"image" Validate:"required"`
	IsFactory    bool   `json:"isFactory,omitempty"`
	Created      string `json:"created,omitempty"`
	CreatedBy    string `json:"createdBy,omitempty"`
}

type ReferenceComputeEnvironment struct {
	Reference  string `json:"reference"`
	Identifier string `json:"identifier,omitempty"`
}

type SourceCode struct {
	URL      string `json:"url" Validate:"required"`
	Branch   string `json:"branch,omitempty"`
	CommitID string `json:"commitID,omitempty"`
	SSHKey   string `json:"sshKey,omitempty"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
	if err := s.Validate.Struct(env); err != nil {
		return nil, &Response{}, err
	}
	req, err := s.Client.NewAIRequest("POST", s.path("ComputeEnvironment"), env, nil)
	if err != nil {
		return nil, &Response{}, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var createdEnv ComputeEnvironment
	resp, err := s.Client.Do(req, &createdEnv)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateComputeEnvironment: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdEnv, resp, nil
}

func (s *ComputeEnvironmentService) DeleteComputeEnvironment(env ComputeEnvironment) (*Response, error) {
	req, err := s.Client.NewAIRequest("DELETE", s.path("ComputeEnvironment", env.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	resp, err := s.Client.Do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteComputeEnvironment: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *ComputeEnvironmentService) GetComputeEnvironmentByID(id string) (*ComputeEnvironment, *Response, error) {
	req, err := s.Client.NewAIRequest("GET", s.path("ComputeEnvironment", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var foundEnv ComputeEnvironment
	resp, err := s.Client.Do(req, &foundEnv)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetComputeEnvironmentByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundEnv, resp, nil
}

func (s *ComputeEnvironmentService) GetComputeEnvironments(opt *GetOptions, options ...OptionFunc) ([]ComputeEnvironment, *Response, error) {
	req, err := s.Client.NewAIRequest("GET", s.path("ComputeEnvironment"), opt, options...)
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
	resp, err := s.Client.Do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
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
