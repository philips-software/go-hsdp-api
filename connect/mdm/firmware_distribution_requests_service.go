package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type FirmwareDistributionRequestsService struct {
	*Client
	validate *validator.Validate
}

var (
	firmwareDistributionRequestAPIVersion = "1"
)

// GetFirmwareDistributionRequestOptions struct describes search criteria for looking up FirmwareDistributionRequest
type GetFirmwareDistributionRequestOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

type FirmwareDistributionRequest struct {
	ResourceType              string      `json:"resourceType" validate:"required"`
	ID                        string      `json:"id,omitempty"`
	Status                    string      `json:"status" validate:"required"`
	UserConsentRequired       bool        `json:"userConsentRequired"`
	DistributionTargets       []Reference `json:"distributionTarget" validate:"required,min=1,max=10"`
	FirmwareVersion           string      `json:"firmwareVersion" validate:"required"`
	OrchestrationMode         string      `json:"orchestrationMode" validate:"required,oneof=none continuous snapshot"`
	FirmwareComponentVersions []Reference `json:"firmwareComponentVersions" validate:"required,min=1,max=5"`
	Description               string      `json:"description" validate:"omitempty,max=250"`
}

// Create creates a FirmwareDistributionRequest
func (c *FirmwareDistributionRequestsService) Create(ac FirmwareDistributionRequest) (*FirmwareDistributionRequest, *Response, error) {
	ac.ResourceType = "FirmwareDistributionRequest"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/FirmwareDistributionRequest", ac, nil)
	req.Header.Set("api-version", firmwareDistributionRequestAPIVersion)

	var created FirmwareDistributionRequest

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

// Delete deletes the given FirmwareDistributionRequest
func (c *FirmwareDistributionRequestsService) Delete(ac FirmwareDistributionRequest) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/FirmwareDistributionRequest/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", firmwareDistributionRequestAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *FirmwareDistributionRequestsService) GetByID(id string) (*FirmwareDistributionRequest, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/FirmwareDistributionRequest/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource FirmwareDistributionRequest

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

// Find looks up services based on GetFirmwareDistributionRequestOptions
func (c *FirmwareDistributionRequestsService) Find(opt *GetFirmwareDistributionRequestOptions, options ...OptionFunc) (*[]FirmwareDistributionRequest, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/FirmwareDistributionRequest", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareDistributionRequestAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []FirmwareDistributionRequest
	for _, c := range bundleResponse.Entry {
		var resource FirmwareDistributionRequest
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *FirmwareDistributionRequestsService) Update(ac FirmwareDistributionRequest) (*FirmwareDistributionRequest, *Response, error) {
	ac.ResourceType = "FirmwareDistributionRequest"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/FirmwareDistributionRequest/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareDistributionRequestAPIVersion)

	var updated FirmwareDistributionRequest

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
