package dicom

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/philips-software/go-hsdp-api/internal"
)

type Notification struct {
	ID                    string `json:"id,omitempty"`
	Enabled               bool   `json:"enabled" validate:"required"`
	Endpoint              string `json:"endPoint" validate:"required"`
	DefaultOrganizationID string `json:"defaultOrganizationID" validate:"required"`
}

// CreateNotification creates a notification
func (c *ConfigService) CreateNotification(repo Notification, opt *QueryOptions, options ...OptionFunc) (*Notification, *Response, error) {
	bodyBytes, err := json.Marshal(repo)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.newDICOMRequest("POST", "config/dicom/"+c.profile+"/notification", bodyBytes, opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var createdRepo Notification
	resp, err := c.client.do(req, &createdRepo)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateNotification: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdRepo, resp, nil
}

// GetNotifications search for notifications
func (c *ConfigService) GetNotifications(opt *QueryOptions, options ...OptionFunc) (*[]Notification, *Response, error) {
	bodyBytes := []byte("")
	req, err := c.client.newDICOMRequest("GET", "config/dicom/"+c.profile+"/notification", bodyBytes, nil, options...)
	if err != nil {
		return nil, nil, err
	}
	if opt != nil && opt.OrganizationID != nil {
		req.Header.Set("OrganizationID", *opt.OrganizationID)
	}
	req.Header.Set("Content-Type", "application/json")
	var bundleResponse internal.Bundle
	resp, err := c.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var notifications []Notification
	for _, c := range bundleResponse.Entry {
		var resource Notification
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			notifications = append(notifications, resource)
		}
	}
	return &notifications, resp, nil
}
