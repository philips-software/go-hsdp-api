package dicom

import (
	"encoding/json"
	"io"
)

// AdvancedSettings
type AdvancedSettings struct {
	PDULength              int `json:"pduLength"`
	ArtimTimeout           int `json:"artimTimeout"`
	AssociationIdleTimeout int `json:"associationIdleTimeout"`
}

// UnsecuredNetworkConnection
type UnsecuredNetworkConnection struct {
	Port             int              `json:"port"`
	AdvancedSettings AdvancedSettings `json:"advancedSettings"`
}

// ApplicationEntity
type ApplicationEntity struct {
	AllowAny           bool   `json:"allowAny"`
	AeTitle            string `json:"aeTitle"`
	OrganizationID     string `json:"organizationId"`
	AdditionalSettings struct {
		serviceTimeout int `json:"serviceTimeout"`
	} `json:"additionalSettings"`
}

// StoreService
type StoreService struct {
	ID                        string                     `json:"id,omitempty"`
	Title                     string                     `json:"title"`
	UnSecureNetworkConnection UnsecuredNetworkConnection `json:"unSecureNetworkConnection"`
	SecureNetworkConnection   json.RawMessage            `json:"secureNetworkConnection"`
	ApplicationEntities       []ApplicationEntity        `json:"applicationEntities"`
}

// SetStoreService
func (c *ConfigService) SetStoreService(svc StoreService, options ...OptionFunc) (*StoreService, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService StoreService
	resp, err := c.client.do(req, &createdService)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}

	return &createdService, resp, nil
}

// GetStoreService
func (c *ConfigService) GetStoreService(options ...OptionFunc) (*StoreService, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/fhirStore", bodyBytes, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service StoreService
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	return &service, resp, nil
}
