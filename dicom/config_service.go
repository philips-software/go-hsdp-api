package dicom

import (
	"encoding/json"
	"github.com/google/fhir/go/jsonformat"
	"io"
)

// ConfigService
type ConfigService struct {
	client   *Client
	timeZone string
	profile  string
	ma       *jsonformat.Marshaller
	um       *jsonformat.Unmarshaller
}

// CDRServiceAccount
type CDRServiceAccount struct {
	ID         string `json:"id,omitempty"`
	ServiceID  string `json:"serviceId"`
	PrivateKey string `json:"privateKey"`
}

// FHIRStore
type FHIRStore struct {
	ID          string `json:"id,omitempty"`
	MPIEndpoint string `json:"mpiEndPoint"`
}

// CredsServiceAccess
type CredsServiceAccess struct {
	ProductKey     string `json:"productKey"`
	BucketName     string `json:"bucketName"`
	FolderPath     string `json:"folderPath"`
	ServiceAccount struct {
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
	Description       string              `json:"description"`
	AccessType        string              `json:"accessType"`
	CredServiceAccess *CredsServiceAccess `json:"credServiceAccess,omitempty"`
	StaticAccess      *StaticAccess       `json:"staticAccess,omitempty"`
}

// GetOptions describes the fields on which you can search for Groups
type GetOptions struct {
	OrganizationID *string `url:"organizationId,omitempty"`
}

// SetCDRServiceAccount
func (c *ConfigService) SetCDRServiceAccount(svc CDRServiceAccount, options ...OptionFunc) (*CDRServiceAccount, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/cdrServiceAccount", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdCDRServiceAccount CDRServiceAccount
	resp, err := c.client.do(req, &createdCDRServiceAccount)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}

	return &createdCDRServiceAccount, resp, nil
}

// GetCDRServiceAccount
func (c *ConfigService) GetCDRServiceAccount(options ...OptionFunc) (*CDRServiceAccount, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/cdrServiceAccount", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var cdrServiceAccount CDRServiceAccount
	resp, err := c.client.do(req, &cdrServiceAccount)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &cdrServiceAccount, resp, nil
}

// SetFHIRStore
func (c *ConfigService) SetFHIRStore(svc FHIRStore, options ...OptionFunc) (*FHIRStore, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var fhirStore FHIRStore
	resp, err := c.client.do(req, &fhirStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}

	return &fhirStore, resp, nil
}

// GetFHIRStore
func (c *ConfigService) GetFHIRStore(options ...OptionFunc) (*FHIRStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var fhirStore FHIRStore
	resp, err := c.client.do(req, &fhirStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &fhirStore, resp, nil
}
