package dicom

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Repository struct {
	ID                  string `json:"id,omitempty"`
	OrganizationID      string `json:"organizationId"`
	ActiveObjectStoreID string `json:"activeObjectStoreId"`
}

// CreateRepository
func (c *ConfigService) CreateRepository(repo Repository, options ...OptionFunc) (*Repository, *Response, error) {
	bodyBytes, err := json.Marshal(repo)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/dicomRepositories", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdRepo Repository
	resp, err := c.client.do(req, &createdRepo)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &createdRepo, resp, nil
}

// GetRepositories
func (c *ConfigService) GetRepositories(opt *GetOptions, options ...OptionFunc) (*[]Repository, *Response, error) {
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
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &repos, resp, nil
}

// GetObjectStore
func (c *ConfigService) GetRepository(id string, options ...OptionFunc) (*Repository, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/dicomRepositories/"+id, bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var repo Repository
	resp, err := c.client.do(req, &repo)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &repo, resp, nil
}

// DeleteObjectStore
func (c *ConfigService) DeleteRepository(repo Repository, options ...OptionFunc) (bool, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("DELETE", "config/dicom/"+c.profile+"/dicomRepositories/"+repo.ID, bodyBytes, nil, options...)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var res bytes.Buffer
	resp, err := c.client.do(req, &res)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}
