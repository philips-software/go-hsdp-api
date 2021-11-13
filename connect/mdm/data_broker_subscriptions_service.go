package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type DataBrokerSubscriptionsService struct {
	*Client
	validate *validator.Validate
}

var (
	dataBrokerSubscriptionAPIVersion = "1"
)

type DataBrokerSubscription struct {
	ResourceType           string          `json:"resourceType" validate:"required"`
	ID                     string          `json:"id,omitempty"`
	Name                   string          `json:"name" validate:"required,max=20"`
	Description            string          `json:"description" validate:"omitempty,max=250"`
	ServiceAgentId         string          `json:"serviceAgentId" validate:"required"`
	DataSubscriberId       Reference       `json:"dataSubscriberId" validate:"required"`
	DataAdapterId          Reference       `json:"dataAdapterId" validate:"required"`
	AuthenticationMethodId Reference       `json:"authenticationMethodId" validate:"required"`
	DataTypeID             Reference       `json:"dataTypeId" validate:"required"`
	Configuration          json.RawMessage `json:"configuration,omitempty"`
	Meta                   *Meta           `json:"meta,omitempty"`
}

// GetDataBrokerSubscriptionOptions struct describes search criteria for looking up DataBrokerSubscription
type GetDataBrokerSubscriptionOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	PropositionID *string `url:"propositionId,omitempty"`
}

// Create creates a DataBrokerSubscription
func (c *DataBrokerSubscriptionsService) Create(ac DataBrokerSubscription) (*DataBrokerSubscription, *Response, error) {
	ac.ResourceType = "DataBrokerSubscription"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/DataBrokerSubscription", ac, nil)
	req.Header.Set("api-version", dataBrokerSubscriptionAPIVersion)

	var created DataBrokerSubscription

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
func (c *DataBrokerSubscriptionsService) Delete(ac DataBrokerSubscription) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/DataBrokerSubscription/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", dataBrokerSubscriptionAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *DataBrokerSubscriptionsService) GetByID(id string) (*DataBrokerSubscription, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/DataBrokerSubscription/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataBrokerSubscriptionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource DataBrokerSubscription

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
func (c *DataBrokerSubscriptionsService) Find(opt *GetDataBrokerSubscriptionOptions, options ...OptionFunc) (*[]DataBrokerSubscription, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/DataBrokerSubscription", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataBrokerSubscriptionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []DataBrokerSubscription
	for _, c := range bundleResponse.Entry {
		var resource DataBrokerSubscription
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *DataBrokerSubscriptionsService) Update(ac DataBrokerSubscription) (*DataBrokerSubscription, *Response, error) {
	ac.ResourceType = "DataBrokerSubscription"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/DataBrokerSubscription/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", dataBrokerSubscriptionAPIVersion)

	var updated DataBrokerSubscription

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
