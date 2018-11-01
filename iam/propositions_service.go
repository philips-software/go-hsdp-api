package iam

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/jeffail/gabs"
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
		return errMissingName
	}
	if p.OrganizationID == "" {
		return errMissingOrganization
	}
	if p.GlobalReferenceID == "" {
		return errMissingGlobalReference
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
		return nil, resp, errEmptyResults
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

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	props, err := p.parseFromBundle(bundleResponse.Bytes())
	return props, resp, err
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
	count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Proposition/%s", &id)
	if count == 0 {
		return nil, resp, errCouldNoReadResourceAfterCreate
	}
	return p.GetPropositionByID(id)
}

func (p *PropositionsService) parseFromBundle(bundle []byte) (*[]Proposition, error) {
	jsonParsed, err := gabs.ParseJSON(bundle)
	if err != nil {
		return nil, err
	}
	count, ok := jsonParsed.S("total").Data().(float64)
	if !ok || count == 0 {
		return nil, errors.New("empty result")
	}
	propositions := make([]Proposition, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var p Proposition
		p.ID, _ = r.Path("id").Data().(string)
		p.Name, _ = r.Path("name").Data().(string)
		p.Description, _ = r.Path("description").Data().(string)
		p.OrganizationID, _ = r.Path("organizationId").Data().(string)
		p.GlobalReferenceID, _ = r.Path("globalReferenceId").Data().(string)
		// TODO: add meta part
		propositions[i] = p
	}
	return &propositions, nil
}
