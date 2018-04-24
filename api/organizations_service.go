package api

import (
	"github.com/hsdp/go-hsdp-iam/iam"
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

func (o *OrganizationsService) GetOrganizationByID(id string) (*iam.Organization, *Response, error) {
	return o.GetOrganization(&GetOrganizationOptions{ID: &id}, nil)
}

func (o *OrganizationsService) GetOrganization(opt *GetOrganizationOptions, options ...OptionFunc) (*iam.Organization, *Response, error) {
	req, err := o.client.NewIDMRequest("GET", "authorize/identity/Organization", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", OrganizationAPIVersion)

	var bundleResponse interface{}

	resp, err := o.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var org iam.Organization
	org.ParseFromBundle(bundleResponse)
	return &org, resp, err
}
