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
	propositionAPIVersion = "1"
)

// Proposition represents a MDM Proposition entity
type Proposition struct {
	ID                              string     `json:"id,omitempty"`
	ResourceType                    string     `json:"resourceType"`
	Name                            string     `json:"name"`
	Description                     string     `json:"description"`
	OrganizationGuid                Identifier `json:"organizationGuid"`
	GlobalReferenceID               string     `json:"globalReferenceId"`
	DefaultCustomerOrganizationGuid Identifier `json:"defaultCustomerOrganizationGuid"`
	Status                          string     `json:"status,omitempty"`
	ValidationEnabled               bool       `json:"validationEnabled"`
	NotificationEnabled             bool       `json:"notificationEnabled"`
	Meta                            *Meta      `json:"meta,omitempty"`
}

// PropositionsService implements actions on IAM Proposition entities
type PropositionsService struct {
	*Client

	validate *validator.Validate
}

// GetPropositionsOptions specifies what search criteria
// can be used to look for entities
type GetPropositionsOptions struct {
	ID                     *string `url:"_id,omitempty"`
	Count                  *int    `url:"_count,omitempty"`
	Page                   *int    `url:"_page,omitempty"`
	OrganizationID         *string `url:"organizationGuid,omitempty"`
	PropositionID          *string `url:"propositionGuid,omitempty"`
	GlobalReferenceID      *string `url:"globalReferenceGuid,omitempty"`
	Name                   *string `url:"name,omitempty"`
	AuthenticationMethodID *string `url:"authenticationMethodId,omitempty"`
}

// GetPropositionByID retrieves a Proposition by its ID
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

// GetPropositions search for a Proposition entity based on the GetPropositions values
func (p *PropositionsService) GetPropositions(opt *GetPropositionsOptions, options ...OptionFunc) (*[]Proposition, *Response, error) {
	req, err := p.NewRequest(http.MethodGet, "/Proposition", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", propositionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := p.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var props []Proposition
	for _, p := range bundleResponse.Entry {
		var prop Proposition
		if err := json.Unmarshal(p.Resource, &prop); err == nil {
			props = append(props, prop)
		}
	}
	return &props, resp, err
}

// CreateProposition creates a Proposition
func (p *PropositionsService) CreateProposition(prop Proposition) (*Proposition, *Response, error) {
	if err := p.validate.Struct(prop); err != nil {
		return nil, nil, err
	}
	req, err := p.NewRequest(http.MethodPost, "/Proposition", &prop, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var created Proposition

	resp, err := p.Do(req, &created)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}
	ok := resp != nil && resp.StatusCode == http.StatusCreated
	if !ok {
		return nil, resp, fmt.Errorf("CreateProposition failed: resp=%v", resp)
	}
	return p.GetPropositionByID(created.ID)
}
