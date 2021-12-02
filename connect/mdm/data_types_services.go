package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type DataTypesService struct {
	*Client
	validate *validator.Validate
}

var (
	dataTypesAPIVersion = "1"
)

type DataType struct {
	ResourceType  string    `json:"resourceType" validate:"required"`
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description"`
	Tags          []string  `json:"tags,omitempty" validate:"omitempty"`
	PropositionId Reference `json:"propositionId" validate:"required"`
	Meta          *Meta     `json:"meta,omitempty"`
}

// GetDataTypeOptions struct describes search criteria for looking up DataType
type GetDataTypeOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	PropositionID *string `url:"propositionId,omitempty"`
}

// Create creates a DataType
func (c *DataTypesService) Create(ac DataType) (*DataType, *Response, error) {
	ac.ResourceType = "DataType"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/DataType", ac, nil)
	req.Header.Set("api-version", dataTypesAPIVersion)

	var created DataType

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
func (c *DataTypesService) Delete(ac DataType) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/DataType/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", dataTypesAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *DataTypesService) GetByID(id string) (*DataType, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/DataType/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataTypesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource DataType

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
func (c *DataTypesService) Find(opt *GetDataTypeOptions, options ...OptionFunc) (*[]DataType, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/DataType", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataTypesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []DataType
	for _, c := range bundleResponse.Entry {
		var resource DataType
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *DataTypesService) Update(ac DataType) (*DataType, *Response, error) {
	ac.ResourceType = "DataType"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/DataType/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataTypesAPIVersion)

	var updated DataType

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
