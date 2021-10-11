package dicom

import (
	"encoding/json"
	"fmt"
	"io"
)

type ImportService struct {
	ID          string `json:"id,omitempty"`
	AETitle     string `json:"aeTitle"`
	Description string `json:"description,omitempty"`
}

// SetImportService
func (c *ConfigService) SetImportService(svc ImportService, opt *QueryOptions, options ...OptionFunc) (*ImportService, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/importService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService ImportService
	resp, err := c.client.do(req, &createdService)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("SetImportService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdService, resp, nil
}

// GetImportService
func (c *ConfigService) GetImportService(opt *QueryOptions, options ...OptionFunc) (*ImportService, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/importService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service ImportService
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetImportService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &service, resp, nil
}
