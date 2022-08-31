package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type BlobSubscriptionsService struct {
	*Client
	validate *validator.Validate
}

var (
	blobSubscriptionPIVersion = "1"
)

type BlobSubscription struct {
	ResourceType          string     `json:"resourceType" validate:"required"`
	ID                    string     `json:"id,omitempty"`
	Name                  string     `json:"name" validate:"required,max=64"`
	Description           string     `json:"description" validate:"omitempty,max=250"`
	DataTypeId            Reference  `json:"dataTypeId" validate:"required"`
	NotificationTopicGuid Identifier `json:"notificationTopicGuid" validate:"required,dive"`
	Meta                  *Meta      `json:"meta,omitempty"`
}

// GetBlobSubscriptionOptions struct describes search criteria for looking up BlobSubscription
type GetBlobSubscriptionOptions struct {
	ID   *string `url:"_id,omitempty"`
	Name *string `url:"name,omitempty"`
}

// Create creates a BlobSubscription
func (c *BlobSubscriptionsService) Create(ac BlobSubscription) (*BlobSubscription, *Response, error) {
	ac.ResourceType = "BlobSubscription"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/BlobSubscription", ac, nil)
	req.Header.Set("api-version", blobSubscriptionPIVersion)

	var created BlobSubscription

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
func (c *BlobSubscriptionsService) Delete(ac BlobSubscription) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/BlobSubscription/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobSubscriptionPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *BlobSubscriptionsService) GetByID(id string) (*BlobSubscription, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/BlobSubscription/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobSubscriptionPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource BlobSubscription

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
func (c *BlobSubscriptionsService) Find(opt *GetBlobSubscriptionOptions, options ...OptionFunc) (*[]BlobSubscription, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/BlobSubscription", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobSubscriptionPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []BlobSubscription
	for _, c := range bundleResponse.Entry {
		var resource BlobSubscription
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *BlobSubscriptionsService) Update(ac BlobSubscription) (*BlobSubscription, *Response, error) {
	ac.ResourceType = "BlobSubscription"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/BlobSubscription/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobSubscriptionPIVersion)

	var updated BlobSubscription

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
