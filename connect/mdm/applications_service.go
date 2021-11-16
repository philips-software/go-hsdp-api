package mdm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

const (
	applicationAPIVersion = "1"
)

// ApplicationsService implements actions on IAM Application entities
type ApplicationsService struct {
	*Client

	validate *validator.Validate
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
	if apps == nil || len(*apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return &(*apps)[0], resp, err
}

// GetApplicationByName retrieves an Application by its Name
func (a *ApplicationsService) GetApplicationByName(name string) (*Application, *Response, error) {
	apps, resp, err := a.GetApplications(&GetApplicationsOptions{Name: &name}, nil)
	if apps == nil || len(*apps) == 0 {
		return nil, resp, ErrNotFound
	}
	return &(*apps)[0], resp, err
}

// GetApplications search for an Applications entity based on the GetApplicationsOptions values
func (a *ApplicationsService) GetApplications(opt *GetApplicationsOptions, options ...OptionFunc) (*[]Application, *Response, error) {
	req, err := a.NewRequest(http.MethodGet, "/Application", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := a.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var apps []Application
	for _, a := range bundleResponse.Entry {
		var app Application
		if err := json.Unmarshal(a.Resource, &app); err == nil {
			apps = append(apps, app)
		}
	}
	return &apps, resp, nil
}

// CreateApplication creates a Application
func (a *ApplicationsService) CreateApplication(app Application) (*Application, *Response, error) {
	app.ResourceType = "Application"
	if err := a.validate.Struct(app); err != nil {
		return nil, nil, err
	}
	req, err := a.NewRequest(http.MethodPost, "/Application", &app, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var created Application

	resp, err := a.Do(req, &created)
	if err == io.EOF { // EOF is not an error in this case
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, fmt.Errorf("CreateApplication: request failed")
	}
	return &created, resp, nil
}

// UpdateApplication creates a Application
func (a *ApplicationsService) UpdateApplication(app Application) (*Application, *Response, error) {
	app.ResourceType = "Application"
	if err := a.validate.Struct(app); err != nil {
		return nil, nil, err
	}
	req, err := a.NewRequest(http.MethodPost, "/Application", &app, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var updated Application

	resp, err := a.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	if err := internal.CheckResponse(resp.Response); err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
