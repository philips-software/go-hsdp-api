package iam

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs"
)

const (
	organizationAPIVersion = "2"
)

// OrganizationsService implements operations on Organization entities
type OrganizationsService struct {
	client *Client
}

// GetOrganizationOptions describes the criteria for looking up Organizations
type GetOrganizationOptions struct {
	ID          *string `url:"_id,omitempty"`
	ParentOrgID *string `url:"parentOrgId,omitempty"`
	Name        *string `url:"name,omitempty"`
}

// CreateOrganization creates a (sub) organization in IAM
func (o *OrganizationsService) CreateOrganization(organization Organization) (*Organization, *Response, error) {
	organization.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:Organization",
	}

	req, err := o.client.NewRequest(IDM, "POST", "authorize/scim/v2/Organizations", &organization, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var newOrg Organization

	resp, err := o.client.Do(req, &newOrg)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, resp, fmt.Errorf("error creating org: %d", resp.StatusCode)
	}
	return &newOrg, resp, err
}

// UpdateOrganization updates the description of the organization.
func (o *OrganizationsService) UpdateOrganization(org Organization) (*Organization, *Response, error) {
	req, err := o.client.NewRequest(IDM, "PUT", "authorize/scim/v2/Organizations/"+org.ID, &org, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", org.Meta.Version)

	var updatedOrg Organization

	resp, err := o.client.Do(req, &updatedOrg)
	if err != nil {
		return nil, resp, err
	}
	return &updatedOrg, resp, err

}

// GetOrganizationByID retrieves an organization by ID
func (o *OrganizationsService) GetOrganizationByID(id string) (*Organization, *Response, error) {
	var foundOrg Organization

	req, err := o.client.NewRequest(IDM, "GET", "authorize/scim/v2/Organizations/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req, &foundOrg)
	if err != nil {
		return nil, resp, err
	}
	return &foundOrg, resp, nil
}

// GetOrganization retrieves an organization based on the GetOrganizationOptions parameters.
// Deprecated: need to switch to SCIM variant
func (o *OrganizationsService) GetOrganization(opt *GetOrganizationOptions, options ...OptionFunc) (*Organization, *Response, error) {
	req, err := o.client.NewRequest(IDM, "GET", "authorize/identity/Organization", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse bytes.Buffer

	resp, err := o.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	organizations, err := o.parseFromBundle(bundleResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	return &(*organizations)[0], resp, nil
}

func (o *OrganizationsService) parseFromBundle(bundle []byte) (*[]Organization, error) {
	jsonParsed, err := gabs.ParseJSON(bundle)
	if err != nil {
		return nil, err
	}
	count, ok := jsonParsed.S("total").Data().(float64)
	if !ok || count == 0 {
		return nil, ErrEmptyResults
	}
	organizations := make([]Organization, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var org Organization
		org.ID, _ = r.Path("resource.id").Data().(string)
		org.Name, _ = r.Path("resource.name").Data().(string)
		org.Description, _ = r.Path("resource.text").Data().(string)
		organizations[i] = org
	}
	return &organizations, nil
}
