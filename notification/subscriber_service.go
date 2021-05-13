package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type SubscriberService struct {
	client *Client

	validate *validator.Validate
}

type Subscriber struct {
	ID                            string `json:"_id,omitempty"`
	ResourceType                  string `json:"resourceType,omitempty"`
	ManagingOrganizationID        string `json:"managingOrganizationId" validate:"required"`
	ManagingOrganization          string `json:"managingOrganization,omitempty"`
	SubscriberProductName         string `json:"subscriberProductName" validate:"required"`
	SubscriberServicename         string `json:"subscriberServiceName" validate:"required"`
	SubscriberServiceinstanceName string `json:"subscriberServiceInstanceName,omitempty"`
	SubscriberServiceBaseURL      string `json:"subscriberServiceBaseUrl" validate:"required"`
	SubscriberServicePathURL      string `json:"subscriberServicePathUrl" validate:"required"`
	Description                   string `json:"description,omitempty"`
}

func (p *SubscriberService) CreateSubscriber(subscriber Subscriber) (*Subscriber, *Response, error) {
	if err := p.validate.Struct(subscriber); err != nil {
		return nil, nil, err
	}
	req, err := p.client.newNotificationRequest("POST", "core/notification/Subscriber", subscriber, nil)
	if err != nil {
		return nil, nil, err
	}
	var createdSubscriber Subscriber
	resp, err := p.client.do(req, &createdSubscriber)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateSubscriber: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdSubscriber, resp, nil
}

func (p *SubscriberService) GetSubscribers(opt *GetOptions, options ...OptionFunc) ([]*Subscriber, *Response, error) {
	var subscribers []*Subscriber

	req, err := p.client.newNotificationRequest("GET", "core/notification/Subscriber", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var bundleResponse internal.Bundle

	resp, err := p.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return subscribers, resp, ErrEmptyResult
	}
	for _, e := range bundleResponse.Entry {
		c := new(Subscriber)
		if err := json.Unmarshal(e.Resource, c); err == nil {
			subscribers = append(subscribers, c)
		} else {
			return nil, resp, err
		}
	}
	return subscribers, resp, err
}

func (p *SubscriberService) DeleteSubscriber(subscriber Subscriber) (bool, *Response, error) {
	req, err := p.client.newNotificationRequest("DELETE", "core/notification/Subscriber/"+subscriber.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", APIVersion)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}
