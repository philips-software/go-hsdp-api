package notification

import (
	"fmt"
	"io"
)

type PublishRequest struct {
	TopicID string `json:"topicId"`
	Message string `json:"message"`
}

type PublishResponse struct {
	ID           string `json:"_id,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	TopicID      string `json:"topicId"`
}

func (c *Client) Publish(request PublishRequest) (*PublishResponse, *Response, error) {
	if err := c.validate.Struct(request); err != nil {
		return nil, nil, err
	}
	req, err := c.newNotificationRequest("POST", "core/notification/Publish", request, nil)
	if err != nil {
		return nil, nil, err
	}
	var publishResponse PublishResponse
	resp, err := c.do(req, &publishResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("publish: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &publishResponse, resp, nil
}
