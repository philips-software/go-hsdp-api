package notification

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ProducerService struct {
	client *Client

	validate *validator.Validate
}

type Producer struct {
	ID                          string `json:"_id,omitempty"`
	ResourceType                string `json:"resourceType,omitempty"`
	ManagingOrganizationID      string `json:"managingOrganizationId" validate:"required"`
	ManagingOrganization        string `json:"managingOrganization,omitempty"`
	ProducerProductName         string `json:"producerProductName" validate:"required"`
	ProducerServiceName         string `json:"producerServiceName" validate:"required"`
	ProducerServiceInstanceName string `json:"producerServiceInstanceName,omitempty" validate:"required"`
	ProducerServiceBaseURL      string `json:"producerServiceBaseUrl" validate:"required"`
	ProducerServicePathURL      string `json:"producerServicePathUrl" validate:"required"`
	Description                 string `json:"description,omitempty"`
}

// GetOptions describes the fields on which you can search for producers
type GetOptions struct {
	ID                    *string `url:"_id,omitempty"`
	ManagedOrganizationID *string `url:"managedOrganizationId,omitempty"`
	ManagedOrganization   *string `url:"managedOrganization,omitempty"`
	ProducerProductName   *string `url:"producerProductName,omitempty"`
	ProducerServiceName   *string `url:"producerServiceName,omitempty"`
	Scope                 *string `url:"scope,omitempty"`
}

func (p *ProducerService) CreateProducer(producer Producer) (*Producer, *Response, error) {
	if err := p.validate.Struct(producer); err != nil {
		return nil, nil, err
	}
	req, err := p.client.newNotificationRequest("POST", "core/notification/Producer", producer, nil)
	if err != nil {
		return nil, nil, err
	}
	var createdProducer Producer
	resp, err := p.client.do(req, &createdProducer)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateProducer: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdProducer, resp, nil
}

func (p *ProducerService) GetProducers(opt *GetOptions, options ...OptionFunc) ([]Producer, *Response, error) {
	var producers []Producer

	req, err := p.client.newNotificationRequest("GET", "core/notification/Producer", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var bundleResponse struct {
		ResourceType string     `json:"resourceType,omitempty"`
		Type         string     `json:"type,omitempty"`
		Total        int        `json:"total"`
		Entry        []Producer `json:"entry"`
	}

	resp, err := p.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return producers, resp, ErrEmptyResult
	}
	for _, e := range bundleResponse.Entry {
		producers = append(producers, e)
	}
	return producers, resp, err
}

func (p *ProducerService) GetProducer(id string) (*Producer, *Response, error) {
	producers, resp, err := p.GetProducers(&GetOptions{ID: &id})
	if err != nil {
		return nil, resp, err
	}
	if producers == nil || len(producers) != 1 {
		return nil, resp, fmt.Errorf("GetProducer: not found")
	}
	return &producers[0], resp, nil
}

func (p *ProducerService) DeleteProducer(producer Producer) (bool, *Response, error) {
	req, err := p.client.newNotificationRequest("DELETE", "core/notification/Producer/"+producer.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", APIVersion)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, fmt.Errorf("DeleteProducer: HTTP %d", resp.StatusCode)
	}
	return true, resp, err
}
