package iam

import (
	"bytes"
	"fmt"

	"net/http"

	"github.com/jeffail/gabs"
)

const (
	organizationAPIVersion = "1"
)

// Organization represents a IAM Organization resource
type Organization struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	DistinctName   string `json:"distinctName,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
}

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
func (o *OrganizationsService) CreateOrganization(parentOrgID, name, description string) (*Organization, *Response, error) {
	var newOrg Organization

	newOrg.Name = name
	newOrg.Description = description

	req, err := o.client.NewIDMRequest("POST", "security/organizations/"+parentOrgID+"/childorganizations", &newOrg, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)

	var bundleResponse bytes.Buffer

	resp, err := o.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("error creating org: %d", resp.StatusCode)
	}
	j, err := gabs.ParseJSON(bundleResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	newOrg.Name = j.Path("exchange.name").Data().(string)
	newOrg.Description = j.Path("exchange.description").Data().(string)
	newOrg.OrganizationID = j.Path("exchange.organizationId").Data().(string)
	return &newOrg, resp, err
}

// UpdateOrganization updates the description of the organization.
func (o *OrganizationsService) UpdateOrganization(org Organization) (*Organization, *Response, error) {
	var updateRequest struct {
		Description string `json:"description"`
	}
	updateRequest.Description = org.Description
	req, err := o.client.NewIDMRequest("PUT", "security/organizations/"+org.OrganizationID, &updateRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var responseBody bytes.Buffer

	resp, err := o.client.Do(req, &responseBody)
	if err != nil {
		return nil, resp, err
	}
	return &org, resp, err

}

// GetOrganizationByID retrieves an organization by ID
func (o *OrganizationsService) GetOrganizationByID(id string) (*Organization, *Response, error) {
	return o.GetOrganization(&GetOrganizationOptions{ID: &id}, nil)
}

// GetOrganization retrieves an organization based on the GetOrganizationOptions parameters.
func (o *OrganizationsService) GetOrganization(opt *GetOrganizationOptions, options ...OptionFunc) (*Organization, *Response, error) {
	req, err := o.client.NewIDMRequest("GET", "authorize/identity/Organization", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)

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
		return nil, errEmptyResults
	}
	organizations := make([]Organization, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var org Organization
		org.OrganizationID, _ = r.Path("resource.id").Data().(string)
		org.Name, _ = r.Path("resource.name").Data().(string)
		org.Description, _ = r.Path("resource.text").Data().(string)
		organizations[i] = org
	}
	return &organizations, nil
}
