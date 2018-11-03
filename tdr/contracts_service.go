package tdr

import (
	"bytes"
	"net/http"
)

// ContractsService provides operations on TDR contracts
type ContractsService struct {
	client *Client
}

// Constants
const (
	TDRAPIVersion = "4"
)

// GetContractOptions describes the fileds on which you can search for Groups
type GetContractOptions struct {
	Organization *string `url:"organization,omitempty"`
	Datatype     *string `url:"dataType,omitempty"`
	Count        *int    `url:"_count,omitempty"`
}

// GetContract searches for contracts in TDR
func (c *ContractsService) GetContract(opt *GetContractOptions, options ...OptionFunc) ([]*Contract, *Response, error) {

	req, err := c.client.NewTDRRequest("GET", "store/tdr/Contract", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", TDRAPIVersion)

	var bundleResponse struct {
		Type         string      `json:"type,omitempty"`
		Total        int         `json:"total,omitempty"`
		Entry        []*Contract `json:"entry,omitempty"`
		ResourceType string      `json:"resourceType,omitempty"`
	}

	resp, err := c.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}

	return bundleResponse.Entry, resp, err
}

// CreateContract creates a new contract in TDR
func (c *ContractsService) CreateContract(contract Contract) (bool, *Response, error) {
	req, err := c.client.NewTDRRequest("POST", "store/tdr/Contract", &contract, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Api-Version", TDRAPIVersion)

	var createResponse bytes.Buffer
	resp, err := c.client.Do(req, &createResponse)
	if err != nil {
		return false, resp, err
	}
	if resp.StatusCode != http.StatusCreated {
		return false, resp, err
	}
	// TODO: capture Location header content
	return true, resp, nil
}
