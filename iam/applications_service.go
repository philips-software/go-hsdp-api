package iam

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
func (o *ApplicationsService) GetApplicationByID(id string) (*Application, *Response, error) {
	return o.GetApplication(&GetApplicationsOptions{ID: &id}, nil)
}

// GetApplication search for an Application entity based on the GetApplicationsOptions values
func (o *ApplicationsService) GetApplication(opt *GetApplicationsOptions, options ...OptionFunc) (*Application, *Response, error) {
	req, err := o.client.NewIDMRequest("GET", "authorize/identity/Application", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", applicationAPIVersion)

	var bundleResponse interface{}

	resp, err := o.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var app Application
	app.parseFromBundle(bundleResponse)
	return &app, resp, err
}
