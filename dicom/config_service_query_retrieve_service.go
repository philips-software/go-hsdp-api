package dicom

import (
	"encoding/json"
	"fmt"
	"io"
)

// SetQueryRetrieveService
func (c *ConfigService) SetQueryRetrieveService(svc BrokenSCPConfig, opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/moveService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService BrokenSCPConfig
	resp, err := c.client.do(req, &createdService)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("SetQueryRetrieveService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdService, resp, nil
}

// GetQueryRetrieveService
func (c *ConfigService) GetQueryRetrieveService(opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/moveService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service []BrokenSCPConfig
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetQueryRetrieveService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	if len(service) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &service[0], resp, nil
}
