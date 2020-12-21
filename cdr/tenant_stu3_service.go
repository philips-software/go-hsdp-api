package cdr

import (
	"bytes"
	"io"
	"net/http"

	"github.com/google/fhir/go/jsonformat"

	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

type TenantSTU3Service struct {
	rootOrgID string
	client    *Client
	timeZone  string
	ma        *jsonformat.Marshaller
	um        *jsonformat.Unmarshaller
}

// Onboard onboards the organization on the CDR under the rootOrgID
func (t *TenantSTU3Service) Onboard(organization *stu3pb.Organization, options ...OptionFunc) (*stu3pb.Organization, *Response, error) {
	organizationJSON, err := t.ma.MarshalResource(organization)
	if err != nil {
		return nil, nil, err
	}
	orgID := organization.Identifier[0].GetValue().Value

	req, err := t.client.NewCDRRequest(http.MethodPut, "store/fhir/"+t.rootOrgID+"/Organization/"+orgID, organizationJSON, options)
	if err != nil {
		return nil, nil, err
	}
	var onboardResponse bytes.Buffer
	resp, err := t.client.Do(req, &onboardResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	unmarshalled, err := t.um.Unmarshal(onboardResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	onboardedOrg := contained.GetOrganization()
	return onboardedOrg, resp, nil
}

func (t *TenantSTU3Service) GetOrganizationByID(orgID string) (*stu3pb.Organization, *Response, error) {
	req, err := t.client.NewCDRRequest(http.MethodGet, "store/fhir/"+t.rootOrgID+"/Organization/"+orgID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var getResponse bytes.Buffer
	resp, err := t.client.Do(req, &getResponse)
	if err != nil && err != io.EOF {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	unmarshalled, err := t.um.Unmarshal(getResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	cdrOrg := contained.GetOrganization()
	return cdrOrg, resp, nil
}
