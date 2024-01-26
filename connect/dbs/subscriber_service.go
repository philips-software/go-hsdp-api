package dbs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type SubscribersService struct {
	*Client
	validate *validator.Validate
}

var (
	subscriberAPIVersion = "1"
)

type SQSSubscriberConfig struct {
	ResourceType                  string `json:"resourceType" validate:"required" enum:"SQSSubscriberConfig"`
	NameInfix                     string `json:"nameInfix" validate:"required"`
	Description                   string `json:"description" validate:"required"`
	QueueType                     string `json:"queueType,omitempty" enum:"Standard|FIFO"`
	DeliveryDelaySeconds          int    `json:"deliveryDelaySeconds,omitempty"`
	MessageRetentionPeriod        int    `json:"messageRetentionPeriod,omitempty"`
	ReceiveMessageWaitTimeSeconds int    `json:"receiveMessageWaitTimeSeconds,omitempty"`
	ServerSideEncryption          bool   `json:"serverSideEncryption,omitempty"`
}

type SQSSubscriber struct {
	ID                            string `json:"id"`
	Meta                          *Meta  `json:"meta"`
	Name                          string `json:"name" validate:"required"`
	Description                   string `json:"description" validate:"required"`
	Status                        string `json:"status" validate:"required" enum:"Creating|Deleting|Active|Updating|InError"`
	ErrorMessage                  string `json:"errorMessage,omitempty"`
	ResourceType                  string `json:"resourceType" validate:"required" enum:"SQSSubscriber"`
	QueueName                     string `json:"queueName" validate:"required"`
	QueueType                     string `json:"queueType" validate:"required" enum:"Standard|FIFO"`
	DeliveryDelaySeconds          int    `json:"deliveryDelaySeconds" validate:"required"`
	MessageRetentionPeriod        int    `json:"messageRetentionPeriod" validate:"required"`
	ReceiveMessageWaitTimeSeconds int    `json:"receiveMessageWaitTimeSeconds" validate:"required"`
	ServerSideEncryption          bool   `json:"serverSideEncryption" validate:"required"`
}

type GetSQSSubscriberOptions struct {
	ID          *string `url:"_id,omitempty"`
	Name        *string `url:"name,omitempty"`
	LastUpdated *string `url:"_lastUpdated,omitempty"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	VersionID   string    `json:"versionId,omitempty"`
}

type SQSBundle struct {
	Type  string          `json:"type,omitempty"`
	Entry []SQSSubscriber `json:"entry,omitempty"`
}

func (b *SubscribersService) CreateSQS(sqsConfig SQSSubscriberConfig) (*SQSSubscriber, *Response, error) {
	sqsConfig.ResourceType = "SQSSubscriberConfig"
	if err := b.validate.Struct(sqsConfig); err != nil {
		return nil, nil, err
	}

	req, _ := b.NewRequest(http.MethodPost, "/Subscriber/SQS", sqsConfig, nil)
	req.Header.Set("api-version", subscriberAPIVersion)

	var created SQSSubscriber

	resp, err := b.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

func (b *SubscribersService) GetSQSByID(id string) (*SQSSubscriber, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Subscriber/SQS/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", subscriberAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource SQSSubscriber

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetSQSByID: %w", err)
	}
	if resource.ID != id {
		return nil, nil, fmt.Errorf("returned resource does not match")
	}
	return &resource, resp, nil
}

func (b *SubscribersService) FindSQS(opt *GetSQSSubscriberOptions, options ...OptionFunc) (*[]SQSSubscriber, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Subscriber/SQS", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", subscriberAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse SQSBundle

	resp, err := b.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}

	return &bundleResponse.Entry, resp, err
}

func (b *SubscribersService) DeleteSQS(subscriber SQSSubscriber) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/Subscriber/SQS/"+subscriber.ID, nil, nil)
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
