package mdm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

var (
	clientAPIVersion = "1"
)

type OAuthClient struct {
	ResourceType           string      `json:"resourceType" validate:"required"`
	ID                     string      `json:"id,omitempty"`
	Name                   string      `json:"name" validate:"required"`
	Description            string      `json:"description"`
	ApplicationId          Reference   `json:"applicationId" validate:"required"`
	GlobalReferenceID      string      `json:"globalReferenceId" validate:"required"`
	RedirectionURIs        []string    `json:"redirectionURIs"`
	ResponseTypes          []string    `json:"responseTypes"`
	UserClient             bool        `json:"userClient"`
	BootstrapClientGuid    *Identifier `json:"bootstrapClientGuid,omitempty"`
	BootstrapClientID      string      `json:"bootstrapClientId,omitempty"`
	BootstrapClientSecret  string      `json:"bootstrapClientSecret,omitempty"`
	BootstrapClientRevoked bool        `json:"bootstrapClientRevoked,omitempty"`
	ClientGuid             *Identifier `json:"clientGuid,omitempty"`
	ClientID               string      `json:"clientId,omitempty"`
	ClientSecret           string      `json:"clientSecret,omitempty"`
	ClientRevoked          bool        `json:"clientRevoked"`
	Meta                   *Meta       `json:"meta,omitempty"`
}

type Reference struct {
	Reference string `json:"reference"`
}

// OAuthClientsService provides operations on IAM roles resources
type OAuthClientsService struct {
	*Client
	validate *validator.Validate
}

// GetOAuthClientsOptions describes search criteria for looking up roles
type GetOAuthClientsOptions struct {
	ID                *string `url:"_id,omitempty"`
	Name              *string `url:"name,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	ApplicationID     *string `url:"applicationId,omitempty"`
}

// CreateOAuthClient creates a Client
func (c *OAuthClientsService) CreateOAuthClient(ac OAuthClient) (*OAuthClient, *Response, error) {
	ac.ResourceType = "OAuthClient"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/OAuthClient", ac, nil)
	req.Header.Set("api-version", clientAPIVersion)

	var createdClient OAuthClient

	resp, err := c.Do(req, &createdClient)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	if resp == nil {
		return nil, resp, fmt.Errorf("CreateOAuthClient (resp=nil): %w", ErrCouldNoReadResourceAfterCreate)
	}

	return &createdClient, resp, nil
}

// DeleteOAuthClient deletes the given Client
func (c *OAuthClientsService) DeleteOAuthClient(ac OAuthClient) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/OAuthClient/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetOAuthClientByID finds a client by its ID
func (c *OAuthClientsService) GetOAuthClientByID(id string) (*OAuthClient, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/OAuthClient/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var client OAuthClient

	resp, err := c.Do(req, &client)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetOAuthClientByID: %w", err)
	}
	return &client, resp, nil
}

// GetOAuthClients looks up clients based on GetClientsOptions
func (c *OAuthClientsService) GetOAuthClients(opt *GetOAuthClientsOptions, options ...OptionFunc) (*[]OAuthClient, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/OAuthClient", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var clients []OAuthClient
	for _, c := range bundleResponse.Entry {
		var client OAuthClient
		if err := json.Unmarshal(c.Resource, &client); err == nil {
			clients = append(clients, client)
		}
	}
	return &clients, resp, err
}

// UpdateScopes updates a clients scope
func (c *OAuthClientsService) UpdateScopes(ac OAuthClient, scopes []string, defaultScopes []string) (bool, *Response, error) {
	return c.UpdateScopesByFlag(ac, scopes, defaultScopes, false)
}

// UpdateScopes updates a clients scope, with possibility to choose between regular client and bootstrap client
func (c *OAuthClientsService) UpdateScopesByFlag(ac OAuthClient, scopes []string, defaultScopes []string, isBootstrapClient bool) (bool, *Response, error) {
	if isBootstrapClient {
		if ac.BootstrapClientGuid == nil {
			return false, nil, fmt.Errorf("missing required IAM bootstrapClientGuid")
		}
	} else {
		if ac.ClientGuid == nil {
			return false, nil, fmt.Errorf("missing required IAM clientGuid")
		}
	}

	var requestBody = struct {
		Scopes        []string `json:"scopes"`
		DefaultScopes []string `json:"defaultScopes"`
	}{
		scopes,
		defaultScopes,
	}

	var clientGUID = ac.ClientGuid.Value
	if isBootstrapClient {
		clientGUID = ac.BootstrapClientGuid.Value
	}

	req, err := c.NewRequest(http.MethodPut, "/OAuthClient/"+ac.ID+"/Client/"+clientGUID+"/$scopes", requestBody)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var putResponse bytes.Buffer

	resp, err := c.Do(req, &putResponse)
	if err != nil {
		return false, resp, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return false, resp, ErrOperationFailed
	}
	return true, resp, nil
}

// Update updates a client
func (c *OAuthClientsService) Update(ac OAuthClient) (*OAuthClient, *Response, error) {
	ac.ResourceType = "OAuthClient"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/OAuthClient/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var updatedClient OAuthClient

	resp, err := c.Do(req, &updatedClient)
	if err != nil {
		return nil, resp, err
	}
	return &updatedClient, resp, nil
}
