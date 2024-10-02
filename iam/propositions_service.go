package iam

import (
	"bytes"
	"fmt"
	"io"
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

// PropositionStatus holds the status of a delete Proposition operation
type PropositionStatus struct {
	Schemas        []string `json:"schemas"`
	ID             string   `json:"id"`
	Status         string   `json:"status"`
	TotalResources int      `json:"totalResources"`
	Meta           *Meta    `json:"meta"`
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

// GetProposition find a Proposition based on the GetPropositions values
func (p *PropositionsService) GetProposition(opt *GetPropositionsOptions, options ...OptionFunc) (*Proposition, *Response, error) {
	props, resp, err := p.GetPropositions(opt, options...)
	if err != nil {
		return nil, resp, err
	}
	if len(*props) == 0 {
		return nil, resp, fmt.Errorf("GetProposition: %w", ErrEmptyResults)
	}
	return &(*props)[0], resp, nil
}

// GetPropositions search for an Proposition entity based on the GetPropositions values
func (p *PropositionsService) GetPropositions(opt *GetPropositionsOptions, options ...OptionFunc) (*[]Proposition, *Response, error) {
	req, err := p.client.newRequest(IDM, "GET", "authorize/identity/Proposition", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int           `json:"total"`
		Entry []Proposition `json:"entry"`
	}

	resp, err := p.client.do(req, &bundleResponse)
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
	req, err := p.client.newRequest(IDM, "POST", "authorize/identity/Proposition", &prop, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse interface{}

	resp, err := p.client.do(req, &bundleResponse)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}
	ok := resp != nil && resp.StatusCode() == http.StatusCreated
	if !ok {
		return nil, resp, fmt.Errorf("CreateProposition failed: resp=%v", resp)
	}
	var id string
	count, _ := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Proposition/%s", &id)
	if count == 0 {
		return nil, resp, fmt.Errorf("CreateProposition: %w", ErrCouldNoReadResourceAfterCreate)
	}
	return p.GetPropositionByID(id)
}

func (p *PropositionsService) DeleteProposition(prop Proposition) (bool, *Response, error) {
	req, err := p.client.newRequest(IDM, "DELETE", "authorize/scim/v2/Propositions/"+prop.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Method", "DELETE")

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if err != nil {
		return false, resp, err
	}
	return resp.StatusCode() == http.StatusAccepted, resp, nil
}

// DeleteStatus returns the status of a delete operation on an organization
func (p *PropositionsService) DeleteStatus(id string) (*PropositionStatus, *Response, error) {
	req, err := p.client.newRequest(IDM, "GET", "authorize/scim/v2/Propositions/"+id+"/deleteStatus", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var deleteResponse PropositionStatus

	resp, err := p.client.do(req, &deleteResponse)
	return &deleteResponse, resp, err
}
