package workspace

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

type Service struct {
	client *ai.Client

	validate *validator.Validate
}

type Workspace struct {
	ID                      string                    `json:"id,omitempty"`
	ResourceType            string                    `json:"resourceType"`
	Name                    string                    `json:"name" validate:"required"`
	Description             string                    `json:"description,omitempty"`
	ComputeTarget           ai.ReferenceComputeTarget `json:"computeTarget" validate:"required"`
	SourceCode              ai.SourceCode             `json:"sourceCode"`
	Labels                  []string                  `json:"labels,omitempty"`
	Type                    string                    `json:"type,omitempty"`
	AdditionalConfiguration string                    `json:"additionalConfiguration,omitempty"`
	Created                 string                    `json:"created,omitempty"`
	CreatedBy               string                    `json:"createdBy,omitempty"`
	LastUpdated             string                    `json:"lastUpdated,omitempty"`
}

type AccessURL struct {
	URL string `json:"url"`
}

type LogArtefact struct {
	StartupLog []string `json:"startupLog"`
}

func (s *Service) path(components ...string) string {
	return path.Join(components...)
}

func (s *Service) CreateWorkspace(model Workspace) (*Workspace, *ai.Response, error) {
	if err := s.validate.Struct(model); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewAIRequest("POST", s.path("Workspace"), model, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var createdWorkspace Workspace
	resp, err := s.client.Do(req, &createdWorkspace)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateWorkspace: %w", ai.ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdWorkspace, resp, nil
}

func (s *Service) DeleteWorkspace(ws Workspace) (*ai.Response, error) {
	req, err := s.client.NewAIRequest("DELETE", s.path("Workspace", ws.ID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	resp, err := s.client.Do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteWorkspace: %w", ai.ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}

func (s *Service) GetWorkspaceByID(id string) (*Workspace, *ai.Response, error) {
	req, err := s.client.NewAIRequest("GET", s.path("Workspace", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var foundWorkspace Workspace
	resp, err := s.client.Do(req, &foundWorkspace)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetWorkspaceByID: %w", ai.ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundWorkspace, resp, nil
}

func (s *Service) GetWorkspaces(opt *ai.GetOptions, options ...ai.OptionFunc) ([]Workspace, *ai.Response, error) {
	req, err := s.client.NewAIRequest("GET", s.path("Workspace"), opt, options...)
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
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ai.ErrEmptyResult
		}
		return nil, resp, err
	}
	var workspaces []Workspace
	for _, e := range bundleResponse.Entry {
		var model Workspace
		if err := json.Unmarshal(e.Resource, &model); err == nil {
			workspaces = append(workspaces, model)
		}
	}
	return workspaces, resp, err
}

func (s *Service) StartWorkspace(ws Workspace) (*ai.Response, error) {
	req, err := s.client.NewAIRequest("POST", s.path("Workspace", ws.ID, "$start"), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	return s.client.Do(req, nil)
}

func (s *Service) StopWorkspace(ws Workspace) (*ai.Response, error) {
	req, err := s.client.NewAIRequest("POST", s.path("Workspace", ws.ID, "$stop"), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	return s.client.Do(req, nil)
}

func (s *Service) GetWorkspaceAccessURL(ws Workspace) (*AccessURL, *ai.Response, error) {
	req, err := s.client.NewAIRequest("POST", s.path("Workspace", ws.ID, "$accessUrl"), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var accessURL AccessURL
	resp, err := s.client.Do(req, &accessURL)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ai.ErrEmptyResult
		}
		return nil, resp, err
	}
	return &accessURL, resp, err
}

func (s *Service) GetWorkspaceLogs(ws Workspace) (*LogArtefact, *ai.Response, error) {
	req, err := s.client.NewAIRequest("POST", s.path("Workspace", ws.ID, "$logs"), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", ai.APIVersion)

	var artefact LogArtefact
	resp, err := s.client.Do(req, &artefact)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ai.ErrEmptyResult
		}
		return nil, resp, err
	}
	return &artefact, resp, err
}
