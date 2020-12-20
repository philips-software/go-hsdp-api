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
}

func (t *TenantSTU3Service) Onboard(organization *stu3pb.Organization, options ...OptionFunc) (*stu3pb.Organization, *Response, error) {
	ma, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if err != nil {
		return nil, nil, err
	}
	organizationJSON, err := ma.MarshalResource(organization)
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
	if err != nil && err != io.EOF {
		return nil, resp, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	um, err := jsonformat.NewUnmarshaller(t.timeZone, jsonformat.STU3)
	if err != nil {
		return nil, resp, err
	}
	unmarshalled, err := um.Unmarshal(onboardResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	onboardedOrg := contained.GetOrganization()
	return onboardedOrg, resp, nil
}
