package dicom

import (
	"encoding/json"
	"fmt"
	"io"
)

// AdvancedSettings
type AdvancedSettings struct {
	PDULength              int `json:"pduLength"`
	ArtimTimeout           int `json:"artimTimeout"`
	AssociationIdleTimeout int `json:"associationIdleTimeout"`
}

// ApplicationEntity
type ApplicationEntity struct {
	AllowAny           bool               `json:"allowAny"`
	AeTitle            string             `json:"aeTitle"`
	OrganizationID     string             `json:"organizationId"`
	AdditionalSettings AdditionalSettings `json:"additionalSettings"`
}

type AdditionalSettings struct {
	ServiceTimeout int `json:"serviceTimeout"`
}

// SCPConfig
type SCPConfig struct {
	ID                        string              `json:"id,omitempty"`
	Title                     string              `json:"title"`
	Description               string              `json:"description,omitempty"`
	UnSecureNetworkConnection NetworkConnection   `json:"unSecureNetworkConnection,omitempty"`
	SecureNetworkConnection   NetworkConnection   `json:"secureNetworkConnection,omitempty"`
	ApplicationEntities       []ApplicationEntity `json:"applicationEntities,omitempty"`
	// TODO: TransferCapability
}

// SetStoreService
func (c *ConfigService) SetStoreService(svc SCPConfig, options ...OptionFunc) (*SCPConfig, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/storeService", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService SCPConfig
	resp, err := c.client.do(req, &createdService)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("SetStoreService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdService, resp, nil
}

// GetStoreService
func (c *ConfigService) GetStoreService(options ...OptionFunc) (*SCPConfig, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/storeService", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service SCPConfig
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetStoreService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &service, resp, nil
}
