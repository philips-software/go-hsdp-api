package iam

import (
	"bytes"
	"fmt"

	"net/http"

	"github.com/jeffail/gabs"
)

const (
	OrganizationAPIVersion = "1"
)

type OrganizationsService struct {
	client *Client
}

type GetOrganizationOptions struct {
	ID          *string `url:"_id,omitempty"`
	ParentOrgID *string `url:"parentOrgId,omitempty"`
	Name        *string `url:"name,omitempty"`
}

func (o *OrganizationsService) CreateOrganization(parentOrgID, name, description string) (*Organization, *Response, error) {
	var newOrg Organization

	newOrg.Name = name
	newOrg.Description = description

	req, err := o.client.NewIDMRequest("POST", "security/organizations/"+parentOrgID+"/childorganizations", &newOrg, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

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

func (o *OrganizationsService) GetOrganizationByID(id string) (*Organization, *Response, error) {
	return o.GetOrganization(&GetOrganizationOptions{ID: &id}, nil)
}

func (o *OrganizationsService) GetOrganization(opt *GetOrganizationOptions, options ...OptionFunc) (*Organization, *Response, error) {
	req, err := o.client.NewIDMRequest("GET", "authorize/identity/Organization", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", OrganizationAPIVersion)

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
