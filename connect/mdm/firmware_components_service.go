package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type FirmwareComponentsService struct {
	*Client
	validate *validator.Validate
}

var (
	firmwareComponentAPIVersion = "1"
)

type FirmwareComponent struct {
	ResourceType  string    `json:"resourceType" validate:"require"`
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name" validate:"required,max=255"`
	Description   string    `json:"description" validate:"omitempty,max=250"`
	MainComponent bool      `json:"mainComponent"`
	DeviceTypeId  Reference `json:"deviceTypeId" validate:"required,dive"`
	Meta          *Meta     `json:"meta,omitempty"`
}

// GetFirmwareComponentOptions struct describes search criteria for looking up a FirmwareComponent
type GetFirmwareComponentOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

// Create creates a FirmwareComponent
func (c *FirmwareComponentsService) Create(ac FirmwareComponent) (*FirmwareComponent, *Response, error) {
	ac.ResourceType = "FirmwareComponent"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/FirmwareComponent", ac, nil)
	req.Header.Set("api-version", firmwareComponentAPIVersion)

	var created FirmwareComponent

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

// Delete deletes the given FirmwareComponent
func (c *FirmwareComponentsService) Delete(ac FirmwareComponent) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/FirmwareComponent/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", firmwareComponentAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *FirmwareComponentsService) GetByID(id string) (*FirmwareComponent, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/FirmwareComponent/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource FirmwareComponent

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

// Find looks up services based on GetFirmwareComponentOptions
func (c *FirmwareComponentsService) Find(opt *GetFirmwareComponentOptions, options ...OptionFunc) (*[]FirmwareComponent, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/FirmwareComponent", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []FirmwareComponent
	for _, c := range bundleResponse.Entry {
		var resource FirmwareComponent
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *FirmwareComponentsService) Update(ac FirmwareComponent) (*FirmwareComponent, *Response, error) {
	ac.ResourceType = "FirmwareComponent"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/FirmwareComponent/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentAPIVersion)

	var updated FirmwareComponent

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
