package mdm

import (
	"fmt"
	"io"
	"net/http"
)

const (
	applicationAPIVersion = "1"
)

// ApplicationsService implements actions on IAM Application entities
type ApplicationsService struct {
	*Client
}

// GetApplicationsOptions specifies what search criteria
// can be used to look for entities
type GetApplicationsOptions struct {
	ID                *string `url:"_id,omitempty"`
	PropositionID     *string `url:"propositionId,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceGuid,omitempty"`
	Name              *string `url:"name,omitempty"`
	DefaultGroupID    *string `url:"defaultGroupGuid,omitempty"`
}

// GetApplicationByID retrieves an Application by its ID
func (a *ApplicationsService) GetApplicationByID(id string) (*Application, *Response, error) {
	apps, resp, err := a.GetApplications(&GetApplicationsOptions{ID: &id}, nil)
	if len(apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return apps[0], resp, err
}

// GetApplicationByName retrieves an Application by its Name
func (a *ApplicationsService) GetApplicationByName(name string) (*Application, *Response, error) {
	apps, resp, err := a.GetApplications(&GetApplicationsOptions{ID: &name}, nil)
	if len(apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return apps[0], resp, err
}

// GetApplications search for an Applications entity based on the GetApplicationsOptions values
func (a *ApplicationsService) GetApplications(opt *GetApplicationsOptions, options ...OptionFunc) ([]*Application, *Response, error) {
	req, err := a.NewRequest(http.MethodGet, "/Application", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int
		Entry []*Application
	}

	resp, err := a.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrEmptyResults
	}

	return bundleResponse.Entry, resp, nil
}

// CreateApplication creates a Application
func (a *ApplicationsService) CreateApplication(app Application) (*Application, *Response, error) {
	if err := a.validate.Struct(app); err != nil {
		return nil, nil, err
	}
	req, err := a.NewRequest(http.MethodPost, "/Application", &app, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := a.Do(req, &bundleResponse)
	if err == io.EOF { // EOF is not an error in this case
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, fmt.Errorf("CreateApplication: request failed")
	}
	ok := resp.StatusCode == http.StatusCreated
	if !ok {
		return nil, resp, fmt.Errorf("CreateApplication: failed with StatusCode=%d", resp.StatusCode)
	}
	var id string
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/connect/mdm/Application/%s", &id)
	if err != nil {
		return nil, resp, err
	}
	if count == 0 {
		return nil, resp, fmt.Errorf("CreateApplication: %w", ErrCouldNoReadResourceAfterCreate)
	}
	return a.GetApplicationByID(id)
}
