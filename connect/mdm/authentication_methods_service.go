package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type AuthenticationMethodsService struct {
	*Client
	validate *validator.Validate
}

var (
	authenticationMethodAPIVersion = "1"
)

type AuthenticationMethod struct {
	ResourceType     string      `json:"resourceType" validate:"required"`
	ID               string      `json:"id,omitempty"`
	Name             string      `json:"name" validate:"required,min=1,max=20"`
	Description      string      `json:"description"`
	LoginName        string      `json:"loginName" validate:"required,min=1,max=78"`
	Password         string      `json:"password" validate:"required,min=1,max=50"`
	ClientID         string      `json:"clientId" validate:"required,min=1,max=78"`
	ClientSecret     string      `json:"clientSecret" validate:"required,min=1,max=78"`
	AuthURL          string      `json:"authUrl,omitempty"`
	AuthMethod       string      `json:"authMethod,omitempty"`
	APIVersion       string      `json:"apiVersion,omitempty"`
	OrganizationGuid *Identifier `json:"organizationGuid,omitempty"`
	Meta             *Meta       `json:"meta,omitempty"`
}

// GetAuthenticationMethodOptions struct describes search criteria for looking up authentication methods
type GetAuthenticationMethodOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

// Create creates a AuthenticationMethod
func (c *AuthenticationMethodsService) Create(ac AuthenticationMethod) (*AuthenticationMethod, *Response, error) {
	ac.ResourceType = "AuthenticationMethod"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/AuthenticationMethod", ac, nil)
	req.Header.Set("api-version", authenticationMethodAPIVersion)

	var created AuthenticationMethod

	resp, err := c.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

// Delete deletes the given ServiceAction
func (c *AuthenticationMethodsService) Delete(ac AuthenticationMethod) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/AuthenticationMethod/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", authenticationMethodAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *AuthenticationMethodsService) GetByID(id string) (*AuthenticationMethod, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/AuthenticationMethod/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", authenticationMethodAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource AuthenticationMethod

	resp, err := c.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetByID: %w", err)
	}
	return &resource, resp, nil
}

// Find looks up services based on GetServiceActionOptions
func (c *AuthenticationMethodsService) Find(opt *GetAuthenticationMethodOptions, options ...OptionFunc) (*[]AuthenticationMethod, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/AuthenticationMethod", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", authenticationMethodAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []AuthenticationMethod
	for _, c := range bundleResponse.Entry {
		var resource AuthenticationMethod
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *AuthenticationMethodsService) Update(ac AuthenticationMethod) (*AuthenticationMethod, *Response, error) {
	ac.ResourceType = "AuthenticationMethod"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/AuthenticationMethod/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", authenticationMethodAPIVersion)

	var updated AuthenticationMethod

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
