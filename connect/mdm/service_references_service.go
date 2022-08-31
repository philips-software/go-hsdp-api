package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type ServiceReferencesService struct {
	*Client
	validate *validator.Validate
}

var (
	serviceReferenceAPIVersion = "1"
)

type ServiceReference struct {
	ResourceType      string      `json:"resourceType" validate:"required"`
	ID                string      `json:"id,omitempty"`
	Name              string      `json:"name" validate:"required"`
	Description       string      `json:"description"`
	ApplicationID     Reference   `json:"applicationId" validate:"required"`
	StandardServiceID Reference   `json:"standardServiceId" validate:"required"`
	MatchingRule      string      `json:"matchingRule,omitempty"`
	ServiceActionIDs  []Reference `json:"serviceActionIds" validate:"required,min=1,max=20"`
	BootstrapEnabled  bool        `json:"bootstrapEnabled"`
	Meta              *Meta       `json:"meta,omitempty"`
}

// GetServiceReferenceOptions struct describes search criteria for looking up service references
type GetServiceReferenceOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

// Create creates a ServiceReference
func (c *ServiceReferencesService) Create(ac ServiceReference) (*ServiceReference, *Response, error) {
	ac.ResourceType = "ServiceReference"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/ServiceReference", ac, nil)
	req.Header.Set("api-version", serviceReferenceAPIVersion)

	var created ServiceReference

	resp, err := c.Do(req, &created)

	ok := resp != nil && (resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	if resp == nil {
		return nil, resp, fmt.Errorf("create (resp=nil): %w", ErrCouldNoReadResourceAfterCreate)
	}

	return c.GetByID(created.ID)
}

// Delete deletes the given ServiceAction
func (c *ServiceReferencesService) Delete(ac ServiceReference) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/ServiceReference/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", serviceReferenceAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *ServiceReferencesService) GetByID(id string) (*ServiceReference, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/ServiceReference/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceReferenceAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource ServiceReference

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
func (c *ServiceReferencesService) Find(opt *GetServiceReferenceOptions, options ...OptionFunc) (*[]ServiceReference, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/ServiceReference", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceReferenceAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []ServiceReference
	for _, c := range bundleResponse.Entry {
		var resource ServiceReference
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *ServiceReferencesService) Update(ac ServiceReference) (*ServiceReference, *Response, error) {
	ac.ResourceType = "ServiceReference"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/ServiceReference/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", serviceReferenceAPIVersion)

	var updated ServiceReference

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
