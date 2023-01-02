package iam

import (
	"bytes"
	"fmt"
	"io"
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
	apps, resp, err := a.GetApplications(&GetApplicationsOptions{ID: String(id)}, nil)
	if len(apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return apps[0], resp, err
}

// GetApplicationByName retrieves an Application by its Name
func (a *ApplicationsService) GetApplicationByName(name string) (*Application, *Response, error) {
	apps, resp, err := a.GetApplications(&GetApplicationsOptions{ID: String(name)}, nil)
	if len(apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return apps[0], resp, err
}

// GetApplications search for an Applications entity based on the GetApplicationsOptions values
func (a *ApplicationsService) GetApplications(opt *GetApplicationsOptions, options ...OptionFunc) ([]*Application, *Response, error) {
	req, err := a.client.newRequest(IDM, "GET", "authorize/identity/Application", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int
		Entry []*Application
	}

	resp, err := a.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrEmptyResults
	}

	return bundleResponse.Entry, resp, nil
}

// CreateApplication creates an Application
func (a *ApplicationsService) CreateApplication(app Application) (*Application, *Response, error) {
	if err := a.client.validate.Struct(app); err != nil {
		return nil, nil, err
	}
	req, err := a.client.newRequest(IDM, "POST", "authorize/identity/Application", &app, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := a.client.do(req, &bundleResponse)
	if err == io.EOF { // EOF is not an error in this case
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, fmt.Errorf("CreateApplication: request failed")
	}
	ok := resp.StatusCode() == http.StatusCreated
	if !ok {
		return nil, resp, fmt.Errorf("CreateApplication: failed with StatusCode=%d", resp.StatusCode())
	}
	var id string
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Application/%s", &id)
	if err != nil {
		return nil, resp, err
	}
	if count == 0 {
		return nil, resp, fmt.Errorf("CreateApplication: %w", ErrCouldNoReadResourceAfterCreate)
	}
	return a.GetApplicationByID(id)
}

// DeleteApplication deletes an Application
func (a *ApplicationsService) DeleteApplication(app Application) (bool, *Response, error) {
	req, err := a.client.newRequest(IDM, "DELETE", "authorize/scim/v2/Applications/"+app.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Method", "DELETE")

	var deleteResponse bytes.Buffer

	resp, err := a.client.do(req, &deleteResponse)
	if err != nil {
		return false, resp, err
	}
	return resp.StatusCode() == http.StatusAccepted, resp, nil
}

// DeleteStatus returns the status of a delete operation on an organization
func (a *ApplicationsService) DeleteStatus(id string) (*ApplicationStatus, *Response, error) {
	req, err := a.client.newRequest(IDM, http.MethodGet, "authorize/scim/v2/Organizations/"+id+"/deleteStatus", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var deleteResponse ApplicationStatus

	resp, err := a.client.do(req, &deleteResponse)
	return &deleteResponse, resp, err
}
