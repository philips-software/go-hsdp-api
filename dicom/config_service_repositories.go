package dicom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	ID                  string                  `json:"id,omitempty"`
	OrganizationID      string                  `json:"organizationId"`
	ActiveObjectStoreID string                  `json:"activeObjectStoreId"`
	StoreAsComposite    *bool                   `json:"storeAsComposite,omitempty"`
	Notification        *RepositoryNotification `json:"notification,omitempty"`
}

type RepositoryNotification struct {
	Enabled        bool   `json:"enabled"`
	OrganizationID string `json:"organizationId"`
}

// CreateRepository
func (c *ConfigService) CreateRepository(repo Repository, opt *QueryOptions, options ...OptionFunc) (*Repository, *Response, error) {
	bodyBytes, err := json.Marshal(repo)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/dicomRepositories", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdRepo Repository
	resp, err := c.client.do(req, &createdRepo)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateRepository: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdRepo, resp, nil
}

// GetRepositories
func (c *ConfigService) GetRepositories(opt *QueryOptions, options ...OptionFunc) (*[]Repository, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/dicomRepositories", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	if opt != nil && opt.OrganizationID != nil {
		req.Header.Set("OrganizationID", *opt.OrganizationID)
	}
	req.Header.Set("Content-Type", "application/json")
	var repos []Repository
	resp, err := c.client.do(req, &repos)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetRepositories: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &repos, resp, nil
}

// GetRepository
func (c *ConfigService) GetRepository(id string, opt *QueryOptions, options ...OptionFunc) (*Repository, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/dicomRepositories/"+id, bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var repo Repository
	resp, err := c.client.do(req, &repo)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetRepository: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &repo, resp, nil
}

// DeleteRepository
func (c *ConfigService) DeleteRepository(repo Repository, opt *QueryOptions, options ...OptionFunc) (bool, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("DELETE", "config/dicom/"+c.profile+"/dicomRepositories/"+repo.ID, bodyBytes, opt, options...)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var res bytes.Buffer
	resp, err := c.client.do(req, &res)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteRepository: %w", ErrEmptyResult)
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}
