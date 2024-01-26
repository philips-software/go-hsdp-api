package dbs

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
	"net/http"
)

type SubscriptionService struct {
	*Client
	validate *validator.Validate
}

var (
	subscriptionAPIVersion = "1"
)

type TopicSubscriptionConfig struct {
	ResourceType              string `json:"resourceType" validate:"required" enum:"TopicSubscriptionConfig"`
	NameInfix                 string `json:"nameInfix" validate:"required"`
	Description               string `json:"description" validate:"required"`
	SubscriberId              string `json:"subscriberId" validate:"required"`
	DeliverDataOnly           bool   `json:"deliverDataOnly,omitempty"`
	KinesisStreamPartitionKey string `json:"kinesisStreamPartitionKey,omitempty"`
	DataType                  string `json:"dataType" validate:"required"`
}

type Subscriber struct {
	ID       string `json:"id" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type TopicSubscription struct {
	ResourceType              string      `json:"resourceType" validate:"required" enum:"TopicSubscription"`
	ID                        string      `json:"id"`
	Meta                      *Meta       `json:"meta"`
	Name                      string      `json:"name" validate:"required"`
	Description               string      `json:"description" validate:"required"`
	Subscriber                *Subscriber `json:"subscriber" validate:"required"`
	DeliverDataOnly           bool        `json:"deliverDataOnly,omitempty"`
	KinesisStreamPartitionKey string      `json:"kinesisStreamPartitionKey,omitempty"`
	Status                    string      `json:"status" validate:"required" enum:"Creating|Deleting|Active|Updating|InError"`
	ErrorMessage              string      `json:"errorMessage,omitempty"`
	DataType                  string      `json:"dataType" validate:"required"`
	RuleName                  string      `json:"ruleName" validate:"required"`
}

type GetTopicSubscriptionOptions struct {
	ID          *string `url:"_id,omitempty"`
	Name        *string `url:"name,omitempty"`
	LastUpdated *string `url:"_lastUpdated,omitempty"`
}

type TopicSubscriptionBundle struct {
	Type  string              `json:"type,omitempty"`
	Entry []TopicSubscription `json:"entry,omitempty"`
}

func (b *SubscriptionService) CreateTopicSubscription(subscriptionConfig TopicSubscriptionConfig) (*TopicSubscription, *Response, error) {
	subscriptionConfig.ResourceType = "TopicSubscriptionConfig"
	if err := b.validate.Struct(subscriptionConfig); err != nil {
		return nil, nil, err
	}

	req, _ := b.NewRequest(http.MethodPost, "/Subscription/Topic", subscriptionConfig, nil)
	req.Header.Set("api-version", subscriptionAPIVersion)

	var created TopicSubscription

	resp, err := b.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

func (b *SubscriptionService) GetTopicSubscriptionByID(id string) (*TopicSubscription, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Subscription/Topic/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", subscriptionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource TopicSubscription

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetTopicSubscriptionByID: %w", err)
	}
	if resource.ID != id {
		return nil, nil, fmt.Errorf("returned resource does not match")
	}
	return &resource, resp, nil
}

func (b *SubscriptionService) FindTopicSubscription(opt *GetTopicSubscriptionOptions, options ...OptionFunc) (*[]TopicSubscription, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Subscription/Topic", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", subscriptionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse TopicSubscriptionBundle

	resp, err := b.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}

	return &bundleResponse.Entry, resp, err
}

func (b *SubscriptionService) DeleteTopicSubscription(subscription TopicSubscription) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/Subscription/Topic/"+subscription.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", subscriberAPIVersion)

	var deleteResponse interface{}

	resp, err := b.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}
