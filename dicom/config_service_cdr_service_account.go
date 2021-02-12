package dicom

import (
	"encoding/json"
	"io"
)

// CDRServiceAccount
type CDRServiceAccount struct {
	ID         string `json:"id,omitempty"`
	ServiceID  string `json:"serviceId"`
	PrivateKey string `json:"privateKey"`
}

// Valid
func (c CDRServiceAccount) Valid() bool {
	return c.PrivateKey != "" && c.ServiceID != ""
}

// SetCDRServiceAccount
func (c *ConfigService) SetCDRServiceAccount(svc CDRServiceAccount, opt *QueryOptions, options ...OptionFunc) (*CDRServiceAccount, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/cdrServiceAccount", bodyBytes, opt, options...)
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
func (c *ConfigService) GetCDRServiceAccount(opt *QueryOptions, options ...OptionFunc) (*CDRServiceAccount, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/cdrServiceAccount", bodyBytes, opt, options...)
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
