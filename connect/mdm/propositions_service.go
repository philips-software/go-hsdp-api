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

// Proposition represents an MDM Proposition entity
type Proposition struct {
	ResourceType                          string      `json:"resourceType" validate:"required"`
	ID                                    string      `json:"id,omitempty"`
	Name                                  string      `json:"name" validate:"required"`
	Description                           string      `json:"description"`
	OrganizationGuid                      Identifier  `json:"organizationGuid"`
	PropositionGuid                       *Identifier `json:"propositionGuid,omitempty" validate:"omitempty,dive"`
	GlobalReferenceID                     string      `json:"globalReferenceId"`
	Status                                string      `json:"status,omitempty" validate:"omitempty,oneof=DRAFT ACTIVE"`
	DefaultCustomerOrganizationGuid       *Identifier `json:"defaultCustomerOrganizationGuid,omitempty" validate:"omitempty,dive"`
	ValidationEnabled                     bool        `json:"validationEnabled"`
	ExternalProvisionValidationURL        string      `json:"externalProvisionValidationUrl,omitempty"`
	ExternalProvisionValidationApiVersion string      `json:"externalProvisionValidationApiVersion,omitempty"`
	AuthenticationMethodID                *Reference  `json:"authenticationMethodId,omitempty"`
	NotificationEnabled                   bool        `json:"notificationEnabled"`
	TopicGuid                             *Identifier `json:"topicGuid,omitempty" validate:"omitempty,dive"`
	Meta                                  *Meta       `json:"meta,omitempty"`
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
		return nil, resp, fmt.Errorf("GetProposition: %w", ErrEmptyResult)
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
	prop.ResourceType = "Proposition"
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
	return &created, resp, nil
}

// UpdateProposition creates a Proposition
func (p *PropositionsService) UpdateProposition(prop Proposition) (*Proposition, *Response, error) {
	prop.ResourceType = "Proposition"
	if err := p.validate.Struct(prop); err != nil {
		return nil, nil, err
	}
	req, err := p.NewRequest(http.MethodPut, "/Proposition/"+prop.ID, prop, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var updated Proposition

	resp, err := p.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	if err := internal.CheckResponse(resp.Response); err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
