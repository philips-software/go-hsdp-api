package dicom

import (
	"encoding/json"
	"fmt"
	"io"
)

// SetQueryService
func (c *ConfigService) SetQueryService(svc BrokenSCPConfig, opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/queryService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService BrokenSCPConfig
	resp, err := c.client.do(req, &createdService)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("SetQueryService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdService, resp, nil
}

// GetQueryService
func (c *ConfigService) GetQueryService(opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/queryService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service BrokenSCPConfig
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetQueryService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &service, resp, nil
}
