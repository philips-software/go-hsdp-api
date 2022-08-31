package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

var (
	standardServiceAPIVersion = "1"
)

type StandardService struct {
	ResourceType     string       `json:"resourceType" validate:"required"`
	ID               string       `json:"id,omitempty"`
	Name             string       `json:"name" validate:"required"`
	Description      string       `json:"description"`
	Trusted          bool         `json:"trusted"`
	Tags             []string     `json:"tags" validate:"required,min=1,max=1"`
	ServiceUrls      []ServiceURL `json:"serviceUrls" validate:"required,min=1,max=5"`
	OrganizationGuid *Identifier  `json:"organizationGuid,omitempty"`
	Meta             *Meta        `json:"meta,omitempty"`
}

type ServiceURL struct {
	URL                    string     `json:"url"`
	SortOrder              int        `json:"sortOrder"`
	AuthenticationMethodID *Reference `json:"AuthenticationMethodId,omitempty"`
}

// StandardServicesService provides operations on MDM standard service resources
type StandardServicesService struct {
	*Client

	validate *validator.Validate
}

// GetStandardServiceOptions struct { describes search criteria for looking up standard services
type GetStandardServiceOptions struct {
	ID                *string `url:"_id,omitempty"`
	Name              *string `url:"name,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	ApplicationID     *string `url:"applicationId,omitempty"`
}

// CreateStandardService creates a Client
func (c *StandardServicesService) CreateStandardService(ac StandardService) (*StandardService, *Response, error) {
	ac.ResourceType = "StandardService"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/StandardService", ac, nil)
	req.Header.Set("api-version", standardServiceAPIVersion)

	var created StandardService

	resp, err := c.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

// DeleteStandardService deletes the given Client
func (c *StandardServicesService) DeleteStandardService(ac StandardService) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/StandardService/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", standardServiceAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetStandardServiceByID finds a client by its ID
func (c *StandardServicesService) GetStandardServiceByID(id string) (*StandardService, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/StandardService/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", standardServiceAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var service StandardService

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

// GetStandardServices looks up services based on GetStandardServiceOptions
func (c *StandardServicesService) GetStandardServices(opt *GetStandardServiceOptions, options ...OptionFunc) (*[]StandardService, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/StandardService", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", standardServiceAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var services []StandardService
	for _, c := range bundleResponse.Entry {
		var service StandardService
		if err := json.Unmarshal(c.Resource, &service); err == nil {
			services = append(services, service)
		}
	}
	return &services, resp, err
}

// Update updates a standard service
func (c *StandardServicesService) Update(ac StandardService) (*StandardService, *Response, error) {
	ac.ResourceType = "StandardService"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/StandardService/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", standardServiceAPIVersion)

	var updated StandardService

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
