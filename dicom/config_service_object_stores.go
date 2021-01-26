package dicom

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// CreateObjectStore
func (c *ConfigService) CreateObjectStore(store ObjectStore, options ...OptionFunc) (*ObjectStore, *Response, error) {
	bodyBytes, err := json.Marshal(store)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/objectStores", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdObjectStore ObjectStore
	resp, err := c.client.do(req, &createdObjectStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &createdObjectStore, resp, nil
}

// GetObjectStores
func (c *ConfigService) GetObjectStores(opt *GetOptions, options ...OptionFunc) (*[]ObjectStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/objectStores", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	if opt != nil && opt.OrganizationID != nil {
		req.Header.Set("OrganizationID", *opt.OrganizationID)
	}
	req.Header.Set("Content-Type", "application/json")
	var objectStores []ObjectStore
	resp, err := c.client.do(req, &objectStores)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &objectStores, resp, nil
}

// GetObjectStore
func (c *ConfigService) GetObjectStore(id string, options ...OptionFunc) (*ObjectStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/objectStores/"+id, bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var objectStore ObjectStore
	resp, err := c.client.do(req, &objectStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &objectStore, resp, nil
}

// DeleteObjectStore
func (c *ConfigService) DeleteObjectStore(store ObjectStore, options ...OptionFunc) (bool, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("DELETE", "config/dicom/"+c.profile+"/objectStores/"+store.ID, bodyBytes, options...)
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
