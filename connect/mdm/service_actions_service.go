package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type ServiceActionsService struct {
	*Client
	validate *validator.Validate
}

var (
	serviceActionAPIVersion = "1"
)

type ServiceAction struct {
	ResourceType      string     `json:"resourceType"`
	ID                string     `json:"id,omitempty"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	OrganizationGuid  Identifier `json:"organizationGuid"`
	StandardServiceId Reference  `json:"standardServiceId"`
}

// GetServiceActionOptions struct { describes search criteria for looking up service actions
type GetServiceActionOptions struct {
	ID                *string `url:"_id,omitempty"`
	Name              *string `url:"name,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	ApplicationID     *string `url:"applicationId,omitempty"`
}

// Create creates a ServiceAction
func (c *ServiceActionsService) Create(ac ServiceAction) (*ServiceAction, *Response, error) {
	ac.ResourceType = "ServiceAction"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/ServiceAction", ac, nil)
	req.Header.Set("api-version", serviceActionAPIVersion)

	var created ServiceAction

	resp, err := c.Do(req, &created)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	if resp == nil {
		return nil, resp, fmt.Errorf("create (resp=nil): %w", ErrCouldNoReadResourceAfterCreate)
	}

	return c.GetByID(created.ID)
}

// Delete deletes the given ServiceAction
func (c *ServiceActionsService) Delete(ac ServiceAction) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/ServiceAction/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", serviceActionAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *ServiceActionsService) GetByID(id string) (*ServiceAction, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/ServiceAction/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceActionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var service ServiceAction

	resp, err := c.Do(req, &service)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetStandardServiceByID: %w", err)
	}
	return &service, resp, nil
}

// Find looks up services based on GetServiceActionOptions
func (c *ServiceActionsService) Find(opt *GetServiceActionOptions, options ...OptionFunc) (*[]ServiceAction, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/ServiceAction", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceActionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var services []ServiceAction
	for _, c := range bundleResponse.Entry {
		var service ServiceAction
		if err := json.Unmarshal(c.Resource, &service); err == nil {
			services = append(services, service)
		}
	}
	return &services, resp, err
}

// Update updates a standard service
func (c *ServiceActionsService) Update(ac ServiceAction) (*ServiceAction, *Response, error) {
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/ServiceAction/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceActionAPIVersion)

	var updated ServiceAction

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
