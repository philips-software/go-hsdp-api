package notification

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type TopicService struct {
	client *Client

	validate *validator.Validate
}

type Topic struct {
	ID            string   `json:"_id,omitempty"`
	ResourceType  string   `json:"resourceType,omitempty"`
	Name          string   `json:"name" validate:"required"`
	ProducerID    string   `json:"producerId" validate:"required"`
	Scope         string   `json:"scope" validate:"required"`
	AllowedScopes []string `json:"allowedScopes,omitempty"`
	IsAuditable   bool     `json:"isAuditable,omitempty"`
	Description   string   `json:"description,omitempty"`
}

func (p *TopicService) CreateTopic(topic Topic) (*Topic, *Response, error) {
	if err := p.validate.Struct(topic); err != nil {
		return nil, nil, err
	}
	req, err := p.client.newNotificationRequest("POST", "core/notification/Topic", topic, nil)
	if err != nil {
		return nil, nil, err
	}
	var createdTopic Topic
	resp, err := p.client.do(req, &createdTopic)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateTopic: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdTopic, resp, nil
}

func (p *TopicService) UpdateTopic(topic Topic) (*Topic, *Response, error) {
	if err := p.validate.Struct(topic); err != nil {
		return nil, nil, err
	}
	req, err := p.client.newNotificationRequest("PUT", "core/notification/Topic/"+topic.ID, topic, nil)
	if err != nil {
		return nil, nil, err
	}
	var updateResponse bytes.Buffer
	resp, err := p.client.do(req, &updateResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("UpdateTopic: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return nil, resp, fmt.Errorf("update error: %v", updateResponse)
	}
	updated, resp, err := p.GetTopics(&GetOptions{
		ID: &topic.ID,
	})
	if err != nil {
		return nil, resp, err
	}
	if len(updated) != 1 {
		return nil, resp, fmt.Errorf("failed to retrieve updated Topic %s", topic.ID)
	}
	return &updated[0], resp, nil
}

func (p *TopicService) GetTopics(opt *GetOptions, options ...OptionFunc) ([]Topic, *Response, error) {
	var topics []Topic

	req, err := p.client.newNotificationRequest("GET", "core/notification/Topic", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var bundleResponse struct {
		ResourceType string  `json:"resourceType,omitempty"`
		Type         string  `json:"type,omitempty"`
		Total        int     `json:"total"`
		Entry        []Topic `json:"entry"`
	}

	resp, err := p.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return topics, resp, ErrEmptyResult
	}
	for _, e := range bundleResponse.Entry {
		topics = append(topics, e)
	}
	return topics, resp, err
}

func (p *TopicService) GetTopic(id string) (*Topic, *Response, error) {
	topics, resp, err := p.GetTopics(&GetOptions{ID: &id})
	if err != nil {
		return nil, resp, err
	}
	if topics == nil || len(topics) != 1 {
		return nil, resp, fmt.Errorf("GetTopic: not found")
	}
	return &topics[0], resp, nil
}

func (p *TopicService) DeleteTopic(topic Topic) (bool, *Response, error) {
	req, err := p.client.newNotificationRequest("DELETE", "core/notification/Topic/"+topic.ID, nil, nil)
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
