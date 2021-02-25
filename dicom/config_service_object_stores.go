package dicom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CredsServiceAccess
type CredsServiceAccess struct {
	Endpoint       string `json:"endPoint"`
	ProductKey     string `json:"productKey"`
	BucketName     string `json:"bucketName"`
	FolderPath     string `json:"folderPath"`
	ServiceAccount struct {
		Name                string `json:"name,omitempty"`
		ServiceID           string `json:"serviceId"`
		PrivateKey          string `json:"privateKey"`
		AccessTokenEndPoint string `json:"accessTokenEndPoint"`
		TokenEndPoint       string `json:"tokenEndPoint"`
	} `json:"serviceAccount"`
}

// StaticAccess
type StaticAccess struct {
	Endpoint   string `json:"endPoint"`
	BucketName string `json:"bucketName"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
}

// ObjectStore
type ObjectStore struct {
	ID                string              `json:"id,omitempty"`
	Description       string              `json:"description,omitempty"`
	AccessType        string              `json:"accessType" validate:"required,enum" enum:"direct,s3Creds"`
	CredServiceAccess *CredsServiceAccess `json:"credServiceAccess,omitempty"`
	StaticAccess      *StaticAccess       `json:"staticAccess,omitempty"`
}

// CreateObjectStore
func (c *ConfigService) CreateObjectStore(store ObjectStore, opt *QueryOptions, options ...OptionFunc) (*ObjectStore, *Response, error) {
	bodyBytes, err := json.Marshal(store)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/objectStores", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdObjectStore ObjectStore
	resp, err := c.client.do(req, &createdObjectStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateObjectStore: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdObjectStore, resp, nil
}

// GetObjectStores
func (c *ConfigService) GetObjectStores(opt *QueryOptions, options ...OptionFunc) (*[]ObjectStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/objectStores", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var objectStores []ObjectStore
	resp, err := c.client.do(req, &objectStores)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetObjectStores: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &objectStores, resp, nil
}

// GetObjectStore
func (c *ConfigService) GetObjectStore(id string, opt *QueryOptions, options ...OptionFunc) (*ObjectStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/objectStores/"+id, bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var objectStore ObjectStore
	resp, err := c.client.do(req, &objectStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetObjectStore: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &objectStore, resp, nil
}

// DeleteObjectStore
func (c *ConfigService) DeleteObjectStore(store ObjectStore, opt *QueryOptions, options ...OptionFunc) (bool, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("DELETE", "config/dicom/"+c.profile+"/objectStores/"+store.ID, bodyBytes, opt, options...)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var res bytes.Buffer
	resp, err := c.client.do(req, &res)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteObjectStore: %w", ErrEmptyResult)
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}
