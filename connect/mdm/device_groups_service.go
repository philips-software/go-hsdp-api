package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type DeviceGroupsService struct {
	*Client
	validate *validator.Validate
}

var (
	deviceGroupAPIVersion = "1"
)

type DeviceGroup struct {
	ResourceType     string      `json:"resourceType" validate:"required"`
	ID               string      `json:"id,omitempty"`
	Name             string      `json:"name" validate:"required"`
	Description      string      `json:"description"`
	ApplicationId    Reference   `json:"applicationId"`
	DefaultGroupGuid *Identifier `json:"defaultGroupGuid,omitempty"`
}

// GetDeviceGroupOptions struct describes search criteria for looking up device group
type GetDeviceGroupOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

// Create creates a DeviceGroup
func (c *DeviceGroupsService) Create(ac DeviceGroup) (*DeviceGroup, *Response, error) {
	ac.ResourceType = "DeviceGroup"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/DeviceGroup", ac, nil)
	req.Header.Set("api-version", deviceGroupAPIVersion)

	var created DeviceGroup

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
func (c *DeviceGroupsService) Delete(ac DeviceGroup) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/DeviceGroup/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", deviceGroupAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *DeviceGroupsService) GetByID(id string) (*DeviceGroup, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/DeviceGroup/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceGroupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource DeviceGroup

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
func (c *DeviceGroupsService) Find(opt *GetDeviceGroupOptions, options ...OptionFunc) (*[]DeviceGroup, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/DeviceGroup", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceGroupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []DeviceGroup
	for _, c := range bundleResponse.Entry {
		var resource DeviceGroup
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *DeviceGroupsService) Update(ac DeviceGroup) (*DeviceGroup, *Response, error) {
	ac.ResourceType = "DeviceGroup"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/DeviceGroup/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceGroupAPIVersion)

	var updated DeviceGroup

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
