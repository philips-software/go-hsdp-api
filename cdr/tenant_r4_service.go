package cdr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/google/fhir/go/jsonformat"

	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
)

type TenantR4Service struct {
	client   *Client
	timeZone string
	ma       *jsonformat.Marshaller
	um       *jsonformat.Unmarshaller
}

// Onboard onboards the organization on the CDR under the rootOrgID
func (t *TenantR4Service) Onboard(organization *r4pb.Organization, options ...OptionFunc) (*r4pb.Organization, *Response, error) {
	organizationJSON, err := t.ma.MarshalResource(organization)
	if err != nil {
		return nil, nil, err
	}
	orgID := organization.Identifier[0].GetValue().Value

	req, err := t.client.newCDRRequest(http.MethodPut, fmt.Sprintf("Organization/%s", orgID), organizationJSON, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/fhir+json;fhirVersion=4.0")

	var onboardResponse bytes.Buffer
	resp, err := t.client.do(req, &onboardResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("onboard: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	contained, err := t.um.UnmarshalR4(onboardResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	organization = contained.GetOrganization()
	return organization, resp, nil
}

func (t *TenantR4Service) GetOrganizationByID(orgID string) (*r4pb.Organization, *Response, error) {
	req, err := t.client.newCDRRequest(http.MethodGet, fmt.Sprintf("Organization/%s", orgID), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/fhir+json;fhirVersion=4.0")

	var getResponse bytes.Buffer
	resp, err := t.client.do(req, &getResponse)
	if err != nil && err != io.EOF {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, fmt.Errorf("GetOrganizationByID: %w", ErrEmptyResult)
	}
	contained, err := t.um.UnmarshalR4(getResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	organization := contained.GetOrganization()
	return organization, resp, nil
}
