package iam

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/jeffail/gabs"
)

var (
	clientAPIVersion = "1"
)

// ApplicationClient represents an IAM client resource
type ApplicationClient struct {
	ID                string   `json:"id,omitempty"`
	ClientID          string   `json:"clientId"`
	Type              string   `json:"type"`
	Name              string   `json:"name"`
	Password          string   `json:"password,omitempty"`
	RedirectionURIs   []string `json:"redirectionURIs"`
	ResponseTypes     []string `json:"responseTypes"`
	Description       string   `json:"description"`
	ApplicationID     string   `json:"applicationId"`
	GlobalReferenceID string   `json:"globalReferenceId"`
	IsConsentImplied  bool     `json:"isConsentImplied"`
}

// ClientsService provides operations on IAM roles resources
type ClientsService struct {
	client *Client
}

// GetClientsOptions describes search criteria for looking up roles
type GetClientsOptions struct {
	Name           *string `url:"name,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	RoleID         *string `url:"roleId,omitempty"`
}

// CreateClient creates a Client
func (c *ClientsService) CreateClient(ac ApplicationClient) (*ApplicationClient, *Response, error) {
	req, err := c.client.NewIDMRequest("POST", "authorize/identity/Client", ac, nil)
	req.Header.Set("API-version", clientAPIVersion)

	var createdClient ApplicationClient

	resp, err := c.client.Do(req, &createdClient)
	if err != nil {
		return nil, resp, err
	}
	return &createdClient, resp, err
}

// DeleteClient deletes the given Client
func (c *RolesService) DeleteClient(ac ApplicationClient) (bool, *Response, error) {
	req, err := c.client.NewIDMRequest("DELETE", "authorize/identity/Client/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var deleteResponse interface{}

	resp, err := c.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

// GetClients looks up clients based on GetClientsOptions
func (c *ClientsService) GetClients(opt GetClientsOptions, options ...OptionFunc) (*[]ApplicationClient, *Response, error) {
	req, err := c.client.NewIDMRequest("GET", "authorize/identity/Client", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var bundleResponse bytes.Buffer

	resp, err := c.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	clients, err := c.parseFromBundle(bundleResponse.Bytes())
	return clients, resp, err
}

// UpdateScope updates a clients scope
func (c *ClientsService) UpdateScope(ac ApplicationClient, scopes []string, defaultScope string) (bool, *Response, error) {
	var requestBody = struct {
		Scopes       []string `json:"scopes"`
		DefaultScope string   `json:"defaultScope"`
	}{
		scopes,
		defaultScope,
	}
	req, err := c.client.NewIDMRequest("PUT", "authorize/identity/Client/"+ac.ClientID, requestBody, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var putResponse bytes.Buffer

	resp, err := c.client.DoSigned(req, &putResponse)
	if err != nil {
		return false, resp, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return false, resp, errOperationFailed
	}
	return true, resp, nil
}

func (c *ClientsService) parseFromBundle(bundle []byte) (*[]ApplicationClient, error) {
	jsonParsed, err := gabs.ParseJSON(bundle)
	if err != nil {
		return nil, err
	}
	count, ok := jsonParsed.S("total").Data().(float64)
	if !ok || count == 0 {
		return nil, errors.New("empty result")
	}
	clients := make([]ApplicationClient, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var cl ApplicationClient
		cl.ClientID, _ = r.Path("clientId").Data().(string)
		cl.Name, _ = r.Path("name").Data().(string)
		cl.Description, _ = r.Path("description").Data().(string)
		cl.Type, _ = r.Path("type").Data().(string)
		cl.GlobalReferenceID, _ = r.Path("globalReferenceId").Data().(string)
		cl.ApplicationID, _ = r.Path("applicationId").Data().(string)
		clients[i] = cl
	}
	return &clients, nil
}
