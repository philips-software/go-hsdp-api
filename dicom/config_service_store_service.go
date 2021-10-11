package dicom

import (
	"encoding/json"
	"fmt"
	"io"
)

// AdvancedSettings
type AdvancedSettings struct {
	PDULength              int `json:"pduLength,omitempty"`
	ArtimTimeout           int `json:"artimTimeOut,omitempty"`
	AssociationIdleTimeout int `json:"associationIdleTimeOut,omitempty"`
}

// BrokenAdvancedSettings
type BrokenAdvancedSettings struct {
	PDULength              int `json:"pduLength,omitempty"`
	ArtimTimeout           int `json:"artimTimeout,omitempty"`
	AssociationIdleTimeout int `json:"associationIdleTimeout,omitempty"`
}

// BrokenNetworkConnection
type BrokenNetworkConnection struct {
	Port             int                     `json:"port,omitempty"`
	HostName         string                  `json:"hostName,omitempty"`
	IPAddress        string                  `json:"ipAddress,omitempty"`
	AdvancedSettings *BrokenAdvancedSettings `json:"advancedSettings,omitempty"`
	CertificateInfo  *CertificateInfo        `json:"certificateInfo,omitempty"`
	NetworkTimeout   int                     `json:"networkTimeout,omitempty"`
}

// BrokenSCPConfig
//
// The backstory of the 'Broken' prefix is that at least one JSON field has a difference in casing
// usage e.g.: 'associationIdleTimeout' vs 'associationIdleTimeOut'
// This is a subtle but significant issue i.e. we can't use the regular SCPConfig struct
// in places where the different cased JSON fields are expected. Fixing the API would
// also mean bumping the version number. Until that time we will keep the BrokenSCPConfig
// definition here. We have prefixed it so it stays on the radar.
type BrokenSCPConfig struct {
	ID                        string                   `json:"id,omitempty"`
	Title                     string                   `json:"title"`
	Description               string                   `json:"description,omitempty"`
	UnSecureNetworkConnection *BrokenNetworkConnection `json:"unSecureNetworkConnection,omitempty"`
	SecureNetworkConnection   *BrokenNetworkConnection `json:"secureNetworkConnection,omitempty"`
	ApplicationEntities       []ApplicationEntity      `json:"applicationEntities,omitempty"`
}

// ApplicationEntity
type ApplicationEntity struct {
	AllowAny       bool   `json:"allowAny"`
	AeTitle        string `json:"aeTitle"`
	OrganizationID string `json:"organizationId"`
}

// NetworkConnection
type NetworkConnection struct {
	Port             int               `json:"port,omitempty"`
	HostName         string            `json:"hostName,omitempty"`
	IPAddress        string            `json:"ipAddress,omitempty"`
	DisableIPv6      bool              `json:"disableIpv6"`
	IsSecure         bool              `json:"isSecure"`
	AdvancedSettings *AdvancedSettings `json:"advancedSettings,omitempty"`
	CertificateInfo  *CertificateInfo  `json:"certificateInfo,omitempty"`
	NetworkTimeout   int               `json:"networkTimeout,omitempty"`
}

// SCPConfig
type SCPConfig struct {
	ID                        string              `json:"id,omitempty"`
	Title                     string              `json:"title"`
	Description               string              `json:"description,omitempty"`
	UnSecureNetworkConnection *NetworkConnection  `json:"unSecureNetworkConnection,omitempty"`
	SecureNetworkConnection   *NetworkConnection  `json:"secureNetworkConnection,omitempty"`
	ApplicationEntities       []ApplicationEntity `json:"applicationEntities,omitempty"`
}

// SetStoreService
func (c *ConfigService) SetStoreService(svc BrokenSCPConfig, opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes, err := json.Marshal(svc)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/storeService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdService BrokenSCPConfig
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
func (c *ConfigService) GetStoreService(opt *QueryOptions, options ...OptionFunc) (*BrokenSCPConfig, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/storeService", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var service []BrokenSCPConfig // This will change. The API always return a JSON array now
	resp, err := c.client.do(req, &service)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetStoreService: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	if len(service) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &service[0], resp, nil
}
