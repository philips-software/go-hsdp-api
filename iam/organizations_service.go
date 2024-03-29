package iam

import (
	"bytes"
	"fmt"
	"net/http"
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
	Filter             *string `url:"filter,omitempty"`
	Attributes         *string `url:"attributes,omitempty"`
	ExcludedAttributes *string `url:"excludedAttributes,omitempty"`
}

func FilterOrgEq(orgID string) *GetOrganizationOptions {
	query := "id eq \"" + orgID + "\""
	attributes := "id"
	return &GetOrganizationOptions{
		Filter:     &query,
		Attributes: &attributes,
	}
}

func FilterParentEq(parentID string) *GetOrganizationOptions {
	query := "parent.value eq \"" + parentID + "\""
	attributes := "id"
	return &GetOrganizationOptions{
		Filter:     &query,
		Attributes: &attributes,
	}
}

func FilterNameEq(name string) *GetOrganizationOptions {
	query := "name eq \"" + name + "\""
	attributes := "id"
	return &GetOrganizationOptions{
		Filter:     &query,
		Attributes: &attributes,
	}
}

// CreateOrganization creates a (sub) organization in IAM
func (o *OrganizationsService) CreateOrganization(organization Organization) (*Organization, *Response, error) {
	organization.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:Organization",
	}

	req, err := o.client.newRequest(IDM, http.MethodPost, "authorize/scim/v2/Organizations", &organization, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var newOrg Organization

	resp, err := o.client.do(req, &newOrg)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, resp, fmt.Errorf("error creating org: %d", resp.StatusCode())
	}
	return &newOrg, resp, err
}

// DeleteOrganization deletes the organization
func (o *OrganizationsService) DeleteOrganization(org Organization) (bool, *Response, error) {
	req, err := o.client.newRequest(IDM, "DELETE", "authorize/scim/v2/Organizations/"+org.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Method", "DELETE")

	var deleteResponse bytes.Buffer

	resp, err := o.client.do(req, &deleteResponse)
	if err != nil {
		return false, resp, err
	}
	return resp.StatusCode() == http.StatusAccepted, resp, nil
}

// UpdateOrganization updates the description of the organization.
func (o *OrganizationsService) UpdateOrganization(org Organization) (*Organization, *Response, error) {
	req, err := o.client.newRequest(IDM, "PUT", "authorize/scim/v2/Organizations/"+org.ID, &org, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", org.Meta.Version)

	var updatedOrg Organization

	resp, err := o.client.do(req, &updatedOrg)
	if err != nil {
		return nil, resp, err
	}
	return &updatedOrg, resp, err

}

// GetOrganizationByID retrieves an organization by ID
func (o *OrganizationsService) GetOrganizationByID(id string) (*Organization, *Response, error) {
	var foundOrg Organization

	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Organizations/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.do(req, &foundOrg)
	if err != nil {
		return nil, resp, err
	}
	return &foundOrg, resp, nil
}

// GetOrganization retrieves an organization based on the GetOrganizationOptions parameters.
func (o *OrganizationsService) GetOrganization(opt *GetOrganizationOptions, options ...OptionFunc) (*Organization, *Response, error) {
	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Organizations", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)

	var bundleResponse struct {
		Resources []struct {
			ID string `json:"id"`
		}
	}
	resp, err := o.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if len(bundleResponse.Resources) == 0 {
		return nil, resp, ErrNotFound
	}

	return o.GetOrganizationByID(bundleResponse.Resources[0].ID)
}

// DeleteStatus returns the status of a delete operation on an organization
func (o *OrganizationsService) DeleteStatus(id string) (*OrganizationStatus, *Response, error) {
	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Organizations/"+id+"/deleteStatus", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", organizationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var deleteResponse OrganizationStatus

	resp, err := o.client.do(req, &deleteResponse)
	return &deleteResponse, resp, err
}
