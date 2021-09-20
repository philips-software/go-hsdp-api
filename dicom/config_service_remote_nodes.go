package dicom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NetworkConnection
type NetworkConnection struct {
	Port                          int               `json:"port,omitempty"`
	HostName                      string            `json:"hostName,omitempty"`
	IPAddress                     string            `json:"ipAddress,omitempty"`
	DisableIPv6                   bool              `json:"disableIpv6,omitempty"`
	AdvancedSettings              *AdvancedSettings `json:"advancedSettings,omitempty"`
	CertificateInfo               *CertificateInfo  `json:"certificateInfo,omitempty"`
	AuthenticateClientCertificate bool              `json:"authenticateClientCertificate,omitempty"`
	NetworkTimeout                int               `json:"networkTimeout,omitempty"`
	IsSecure                      bool              `json:"isSecure"`
}

type CertificateInfo struct {
	ID string `json:"id,omitempty"`
}

// RemoteNode
type RemoteNode struct {
	ID                string            `json:"id,omitempty"`
	Title             string            `json:"title"`
	NetworkConnection NetworkConnection `json:"networkConnection"`
	AETitle           string            `json:"aeTitle"`
}

// CreateRemoteNode
func (c *ConfigService) CreateRemoteNode(node RemoteNode, options ...OptionFunc) (*RemoteNode, *Response, error) {
	bodyBytes, err := json.Marshal(node)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/remoteNodes", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdNode RemoteNode
	resp, err := c.client.do(req, &createdNode)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateRemoteNode: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdNode, resp, nil
}

// GetRemoteNodes
func (c *ConfigService) GetRemoteNodes(opt *QueryOptions, options ...OptionFunc) (*[]RemoteNode, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/remoteNodes", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	if opt != nil && opt.OrganizationID != nil {
		req.Header.Set("OrganizationID", *opt.OrganizationID)
	}
	req.Header.Set("Content-Type", "application/json")
	var nodes []RemoteNode
	resp, err := c.client.do(req, &nodes)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetRemoteNodes: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &nodes, resp, nil
}

// GetRemoteNode
func (c *ConfigService) GetRemoteNode(id string, opt *QueryOptions, options ...OptionFunc) (*RemoteNode, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/remoteNodes/"+id, bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var node RemoteNode
	resp, err := c.client.do(req, &node)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetRemoteNode: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &node, resp, nil
}

// DeleteRemoteNode
func (c *ConfigService) DeleteRemoteNode(node RemoteNode, options ...OptionFunc) (bool, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("DELETE", "config/dicom/"+c.profile+"/remoteNodes/"+node.ID, bodyBytes, nil, options...)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var res bytes.Buffer
	resp, err := c.client.do(req, &res)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteRemoteNode: %w", ErrEmptyResult)
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}
