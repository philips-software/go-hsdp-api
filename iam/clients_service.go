package iam

import (
	"bytes"
	"fmt"
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
	Scopes            []string `json:"scopes,omitempty"`
	DefaultScopes     []string `json:"defaultScopes,omitempty"`
	Disabled          bool     `json:"disabled,omitempty"`
	Description       string   `json:"description"`
	ApplicationID     string   `json:"applicationId"`
	GlobalReferenceID string   `json:"globalReferenceId"`
}

// ClientsService provides operations on IAM roles resources
type ClientsService struct {
	client *Client
}

// GetClientsOptions describes search criteria for looking up roles
type GetClientsOptions struct {
	ID                *string `url:"_id,omitempty"`
	Name              *string `url:"name,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	ApplicationID     *string `url:"applicationId,omitempty"`
}

// CreateClient creates a Client
func (c *ClientsService) CreateClient(ac ApplicationClient) (*ApplicationClient, *Response, error) {
	req, err := c.client.NewRequest(IDM, "POST", "authorize/identity/Client", ac, nil)
	req.Header.Set("api-version", clientAPIVersion)

	var createdClient ApplicationClient

	resp, err := c.client.Do(req, &createdClient)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	var id string
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Client/%s", &id)
	if count == 0 {
		return nil, resp, errCouldNoReadResourceAfterCreate
	}
	return c.GetClientByID(id)
}

// DeleteClient deletes the given Client
func (c *ClientsService) DeleteClient(ac ApplicationClient) (bool, *Response, error) {
	req, err := c.client.NewRequest(IDM, "DELETE", "authorize/identity/Client/"+ac.ID, nil, nil)
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

// GetClientByID finds a client by its ID
func (c *ClientsService) GetClientByID(id string) (*ApplicationClient, *Response, error) {
	clients, resp, err := c.GetClients(&GetClientsOptions{ID: &id}, nil)

	if err != nil {
		return nil, resp, err
	}
	if clients == nil {
		return nil, resp, errOperationFailed
	}
	if len(*clients) == 0 {
		return nil, resp, errEmptyResults
	}
	foundClient := (*clients)[0]

	return &foundClient, resp, nil
}

// GetClients looks up clients based on GetClientsOptions
func (c *ClientsService) GetClients(opt *GetClientsOptions, options ...OptionFunc) (*[]ApplicationClient, *Response, error) {
	req, err := c.client.NewRequest(IDM, "GET", "authorize/identity/Client", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse bytes.Buffer

	resp, err := c.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	clients, err := c.parseFromBundle(bundleResponse.Bytes())
	return clients, resp, err
}

// UpdateScope updates a clients scope
func (c *ClientsService) UpdateScopes(ac ApplicationClient, scopes []string, defaultScopes []string) (bool, *Response, error) {
	var requestBody = struct {
		Scopes        []string `json:"scopes"`
		DefaultScopes []string `json:"defaultScopes"`
	}{
		scopes,
		defaultScopes,
	}
	req, err := c.client.NewRequest(IDM, "PUT", "authorize/identity/Client/"+ac.ID+"/$scopes", requestBody, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var putResponse bytes.Buffer

	resp, err := c.client.Do(req, &putResponse)
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
		return nil, errEmptyResults
	}
	clients := make([]ApplicationClient, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var cl ApplicationClient
		cl.ID, _ = r.Path("id").Data().(string)
		cl.ClientID, _ = r.Path("clientId").Data().(string)
		cl.Name, _ = r.Path("name").Data().(string)
		cl.Description, _ = r.Path("description").Data().(string)
		cl.Type, _ = r.Path("type").Data().(string)
		cl.GlobalReferenceID, _ = r.Path("globalReferenceId").Data().(string)
		cl.ApplicationID, _ = r.Path("applicationId").Data().(string)
		children, _ = r.Path("redirectionURIs").Children()
		for _, child := range children {
			cl.RedirectionURIs = append(cl.RedirectionURIs, child.Data().(string))
		}
		children, _ = r.Path("responseTypes").Children()
		for _, child := range children {
			cl.ResponseTypes = append(cl.ResponseTypes, child.Data().(string))
		}
		children, _ = r.Path("scopes").Children()
		for _, child := range children {
			cl.Scopes = append(cl.Scopes, child.Data().(string))
		}
		children, _ = r.Path("defaultScopes").Children()
		for _, child := range children {
			cl.DefaultScopes = append(cl.DefaultScopes, child.Data().(string))
		}
		cl.Disabled, _ = r.Path("disabled").Data().(bool)
		// TODO finish parsing complete resource
		clients[i] = cl
	}
	return &clients, nil
}
