package iam

const (
	propositionAPIVersion = "1"
)

// PropositionsService implements actions on IAM Proposition entities
type PropositionsService struct {
	client *Client
}

// GetPropositionsOptions specifies what search criteria
// can be used to look for entities
type GetPropositionsOptions struct {
	ID                *string `url:"_id,omitempty"`
	PropositionID     *string `url:"propositionId,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	Name              *string `url:"name,omitempty"`
}

// GetPropositionByID retrieves an Proposition by its ID
func (o *PropositionsService) GetPropositionByID(id string) (*Proposition, *Response, error) {
	return o.GetProposition(&GetPropositionsOptions{ID: &id}, nil)
}

// GetProposition search for an Proposition entity based on the GetPropositions values
func (o *PropositionsService) GetProposition(opt *GetPropositionsOptions, options ...OptionFunc) (*Proposition, *Response, error) {
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
	var prop Proposition
	err = prop.parseFromBundle(bundleResponse)
	return &prop, resp, err
}
