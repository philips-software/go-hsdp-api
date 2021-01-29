package dicom

import (
	"encoding/json"
	"io"
)

// SetFHIRStore
func (c *ConfigService) SetFHIRStore(svc FHIRStore, options ...OptionFunc) (*FHIRStore, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, nil, options...)
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
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, nil, options...)
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
