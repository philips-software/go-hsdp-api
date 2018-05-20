package iam

import (
	"fmt"
	"net/http"
)

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
func (p *PropositionsService) GetPropositionByID(id string) (*Proposition, *Response, error) {
	return p.GetProposition(&GetPropositionsOptions{ID: &id}, nil)
}

// GetProposition search for an Proposition entity based on the GetPropositions values
func (p *PropositionsService) GetProposition(opt *GetPropositionsOptions, options ...OptionFunc) (*Proposition, *Response, error) {
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Proposition", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var prop Proposition
	err = prop.parseFromBundle(bundleResponse)
	return &prop, resp, err
}

// CreateProposition creates a Proposition
func (p *PropositionsService) CreateProposition(prop Proposition) (*Proposition, *Response, error) {
	if err := prop.Validate(); err != nil {
		return nil, nil, err
	}
	req, err := p.client.NewIDMRequest("POST", "authorize/identity/Proposition", &prop, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := p.client.Do(req, &bundleResponse)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	var id string
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Proposition/%s", &id)
	if count == 0 {
		return nil, resp, errCouldNoReadResourceAfterCreate
	}
	return p.GetPropositionByID(id)
}

// UpdateProposition updates the Proposition
func (p *PropositionsService) UpdateProposition(prop Proposition) (*Proposition, *Response, error) {
	// TODO: not implemented in v1 as of 2018/05/20
	if true {
		return nil, nil, errNotImplementedByHSDP
	}
	var updateRequest struct {
		Description string `json:"description"`
	}
	updateRequest.Description = prop.Description
	req, err := p.client.NewIDMRequest("PUT", "authorize/identity/Proposition/"+prop.ID, &updateRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var updatedProp Proposition

	resp, err := p.client.Do(req, &updatedProp)
	if err != nil {
		return nil, resp, err
	}
	return &prop, resp, err

}
