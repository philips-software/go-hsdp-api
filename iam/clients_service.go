package iam

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	validator "github.com/go-playground/validator/v10"
)

var (
	clientAPIVersion = "1"
)

// ApplicationClient represents an IAM client resource
type ApplicationClient struct {
	ID                   string   `json:"id,omitempty"`
	ClientID             string   `json:"clientId" validate:"required,min=5,max=20"`
	Type                 string   `json:"type"`
	Name                 string   `json:"name" validate:"required,min=5,max=50"`
	Password             string   `json:"password,omitempty" validate:"required,min=8,max=16"`
	RedirectionURIs      []string `json:"redirectionURIs"`
	ResponseTypes        []string `json:"responseTypes"`
	Scopes               []string `json:"scopes,omitempty"`
	DefaultScopes        []string `json:"defaultScopes,omitempty"`
	Disabled             bool     `json:"disabled,omitempty"`
	Description          string   `json:"description" validate:"max=250"`
	ApplicationID        string   `json:"applicationId" validate:"required"`
	GlobalReferenceID    string   `json:"globalReferenceId" validate:"required,min=3,max=50"`
	ConsentImplied       bool     `json:"consentImplied"`
	AccessTokenLifetime  int      `json:"accessTokenLifetime,omitempty" validate:"min=0,max=31536000"`
	RefreshTokenLifetime int      `json:"refreshTokenLifetime,omitempty" validate:"min=0,max=157680000"`
	IDTokenLifetime      int      `json:"idTokenLifetime,omitempty" validate:"min=0,max=31536000"`
	Meta                 struct {
		VersionID    string    `json:"versionId,omitempty"`
		LastModified time.Time `json:"lastModified,omitempty"`
	} `json:"meta,omitempty"`
}

// ClientsService provides operations on IAM roles resources
type ClientsService struct {
	client *Client

	validate *validator.Validate
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
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	// Remove scopes before calling create
	scopes := ac.Scopes
	defaultScopes := ac.DefaultScopes // Defaults to ["cn"]
	ac.Scopes = []string{}            // Defaults to ["mail", "sn"]
	ac.DefaultScopes = []string{}

	req, _ := c.client.NewRequest(IDM, "POST", "authorize/identity/Client", ac, nil)
	req.Header.Set("api-version", clientAPIVersion)

	var createdClient ApplicationClient

	resp, err := c.client.Do(req, &createdClient)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	var id string
	count, _ := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Client/%s", &id)
	if count == 0 {
		return nil, resp, ErrCouldNoReadResourceAfterCreate
	}
	ac.ID = id
	if len(scopes) > 0 {
		_, _, _ = c.UpdateScopes(ac, scopes, defaultScopes)
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
		return false, resp, err
	}
	return true, resp, nil
}

// GetClientByID finds a client by its ID
func (c *ClientsService) GetClientByID(id string) (*ApplicationClient, *Response, error) {
	clients, resp, err := c.GetClients(&GetClientsOptions{ID: &id}, nil)

	if err != nil {
		return nil, resp, err
	}
	if clients == nil {
		return nil, resp, ErrOperationFailed
	}
	if len(*clients) == 0 {
		return nil, resp, ErrEmptyResults
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

	var bundleResponse struct {
		Total int                 `json:"total"`
		Entry []ApplicationClient `json:"entry"`
	}

	resp, err := c.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &bundleResponse.Entry, resp, err
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
		return false, resp, ErrOperationFailed
	}
	return true, resp, nil
}

// UpdateClient updates a client
func (c *ClientsService) UpdateClient(ac ApplicationClient) (*ApplicationClient, *Response, error) {
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.client.NewRequest(IDM, "PUT", "authorize/identity/Client/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", clientAPIVersion)

	var updatedClient ApplicationClient

	resp, err := c.client.Do(req, &updatedClient)
	if err != nil {
		return nil, resp, err
	}
	return &updatedClient, resp, nil
}
