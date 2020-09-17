package iam

import (
	"fmt"
	"net/http"
)

const (
	propositionAPIVersion = "1"
)

// Proposition represents an IAM Proposition entity
type Proposition struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	OrganizationID    string `json:"organizationId"`
	GlobalReferenceID string `json:"globalReferenceId"`
}

func (p *Proposition) validate() error {
	if p.Name == "" {
		return ErrMissingName
	}
	if p.OrganizationID == "" {
		return ErrMissingOrganization
	}
	if p.GlobalReferenceID == "" {
		return ErrMissingGlobalReference
	}
	return nil
}

// PropositionsService implements actions on IAM Proposition entities
type PropositionsService struct {
	client *Client
}

// GetPropositionsOptions specifies what search criteria
// can be used to look for entities
type GetPropositionsOptions struct {
	ID                *string `url:"_id,omitempty"`
	Count             *int    `url:"_count,omitempty"`
	Page              *int    `url:"_page,omitempty"`
	OrganizationID    *string `url:"organizationId,omitempty"`
	PropositionID     *string `url:"propositionId,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	Name              *string `url:"name,omitempty"`
}

// GetPropositionByID retrieves an Proposition by its ID
func (p *PropositionsService) GetPropositionByID(id string) (*Proposition, *Response, error) {
	return p.GetProposition(&GetPropositionsOptions{ID: &id}, nil)
}

// GetProposition find a Proposition based on the GetPropisitions values
func (p *PropositionsService) GetProposition(opt *GetPropositionsOptions, options ...OptionFunc) (*Proposition, *Response, error) {
	props, resp, err := p.GetPropositions(opt, options...)
	if err != nil {
		return nil, resp, err
	}
	if len(*props) == 0 {
		return nil, resp, ErrEmptyResults
	}
	return &(*props)[0], resp, nil
}

// GetPropositions search for an Proposition entity based on the GetPropositions values
func (p *PropositionsService) GetPropositions(opt *GetPropositionsOptions, options ...OptionFunc) (*[]Proposition, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Proposition", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int           `json:"total"`
		Entry []Proposition `json:"entry"`
	}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &bundleResponse.Entry, resp, err
}

// CreateProposition creates a Proposition
func (p *PropositionsService) CreateProposition(prop Proposition) (*Proposition, *Response, error) {
	if err := prop.validate(); err != nil {
		return nil, nil, err
	}
	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Proposition", &prop, nil)
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
	count, _ := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Proposition/%s", &id)
	if count == 0 {
		return nil, resp, ErrCouldNoReadResourceAfterCreate
	}
	return p.GetPropositionByID(id)
}
