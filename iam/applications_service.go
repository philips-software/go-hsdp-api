package iam

import (
	"fmt"
	"net/http"
)

const (
	applicationAPIVersion = "1"
)

// ApplicationsService implements actions on IAM Application entities
type ApplicationsService struct {
	client *Client
}

// GetApplicationsOptions specifies what search criteria
// can be used to look for entities
type GetApplicationsOptions struct {
	ID                *string `url:"_id,omitempty"`
	PropositionID     *string `url:"propositionId,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	Name              *string `url:"name,omitempty"`
}

// GetApplicationByID retrieves an Application by its ID
func (a *ApplicationsService) GetApplicationByID(id string) (*Application, *Response, error) {
	return a.GetApplication(&GetApplicationsOptions{ID: &id}, nil)
}

// GetApplication search for an Application entity based on the GetApplicationsOptions values
func (a *ApplicationsService) GetApplication(opt *GetApplicationsOptions, options ...OptionFunc) (*Application, *Response, error) {
	req, err := a.client.NewIDMRequest("GET", "authorize/identity/Application", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)

	var bundleResponse interface{}

	resp, err := a.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var app Application
	err = app.parseFromBundle(bundleResponse)
	return &app, resp, err
}

// CreateApplication creates a Application
func (a *ApplicationsService) CreateApplication(app Application) (*Application, *Response, error) {
	if err := app.Validate(); err != nil {
		return nil, nil, err
	}
	req, err := a.client.NewIDMRequest("POST", "authorize/identity/Application", &app, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := a.client.Do(req, &bundleResponse)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	var id string
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Application/%s", &id)
	if count == 0 {
		return nil, resp, errCouldNoReadResourceAfterCreate
	}
	return a.GetApplicationByID(id)
}
