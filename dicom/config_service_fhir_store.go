package dicom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FHIRStore
type FHIRStore struct {
	ID          string `json:"id,omitempty"`
	MPIEndpoint string `json:"mpiEndPoint"`
}

// Valid
func (f FHIRStore) Valid() bool {
	return f.MPIEndpoint != ""
}

// SetFHIRStore
func (c *ConfigService) SetFHIRStore(svc FHIRStore, opt *QueryOptions, options ...OptionFunc) (*FHIRStore, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var fhirStore FHIRStore
	resp, err := c.client.do(req, &fhirStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("SetFHIRStore: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &fhirStore, resp, nil
}

// GetFHIRStore
func (c *ConfigService) GetFHIRStore(opt *QueryOptions, options ...OptionFunc) (*FHIRStore, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var fhirStore FHIRStore
	resp, err := c.client.do(req, &fhirStore)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			return nil, resp, ErrNotFound
		}
		if resp == nil && err != nil {
			err = fmt.Errorf("GetFHIRStore: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &fhirStore, resp, nil
}
