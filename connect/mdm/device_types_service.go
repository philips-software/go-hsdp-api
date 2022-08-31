package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type DeviceTypesService struct {
	*Client
	validate *validator.Validate
}

var (
	deviceTypeAPIVersion = "1"
)

// GetDeviceTypeOptions struct describes search criteria for looking up device types
type GetDeviceTypeOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

type DeviceType struct {
	ResourceType         string          `json:"resourceType" validate:"required"`
	ID                   string          `json:"id,omitempty"`
	Name                 string          `json:"name" validate:"required"`
	Description          string          `json:"description"`
	CTN                  string          `json:"ctn" validate:"required"` // Commercial Type Number
	DeviceGroupId        Reference       `json:"deviceGroupId"`
	DefaultGroupGuid     *Identifier     `json:"defaultGroupGuid,omitempty"`
	CustomTypeAttributes json.RawMessage `json:"customTypeAttributes,omitempty"`
	Meta                 *Meta           `json:"meta,omitempty"`
}

// Create creates a DeviceType
func (c *DeviceTypesService) Create(ac DeviceType) (*DeviceType, *Response, error) {
	ac.ResourceType = "DeviceType"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/DeviceType", ac, nil)
	req.Header.Set("api-version", deviceTypeAPIVersion)

	var created DeviceType

	resp, err := c.Do(req, &created)
	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

// Delete deletes the given DeviceType
func (c *DeviceTypesService) Delete(ac DeviceType) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/DeviceType/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", deviceTypeAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *DeviceTypesService) GetByID(id string) (*DeviceType, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/DeviceType/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceTypeAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource DeviceType

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

// Find looks up services based on GetDeviceTypeOptions
func (c *DeviceTypesService) Find(opt *GetDeviceTypeOptions, options ...OptionFunc) (*[]DeviceType, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/DeviceType", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceTypeAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []DeviceType
	for _, c := range bundleResponse.Entry {
		var resource DeviceType
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *DeviceTypesService) Update(ac DeviceType) (*DeviceType, *Response, error) {
	ac.ResourceType = "DeviceType"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/DeviceType/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceTypeAPIVersion)

	var updated DeviceType

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
