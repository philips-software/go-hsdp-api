package audit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

func (c *Client) CreateAuditEvent(event *dstu2pb.AuditEvent) (*stu3pb.ContainedResource, *Response, error) {
	eventJSON, err := c.ma.MarshalResource(event)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.NewAuditRequest("POST", "core/audit/AuditEvent", eventJSON, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("audit.CreateAuditEvent: %w", err)
	}
	_ = c.httpSigner.SignRequest(req)
	var operationResponse bytes.Buffer
	resp, doErr := c.Do(req, &operationResponse)
	if (doErr != nil && !(doErr == io.EOF || doErr == ErrBadRequest)) || resp == nil {
		if resp == nil && doErr != nil {
			doErr = ErrEmptyResult
		}
		return nil, resp, doErr
	}
	contained := &stu3pb.ContainedResource{}
	if resp.StatusCode == http.StatusCreated {
		return contained, resp, nil
	}
	// OperationOutcome
	unmarshalled, _ := c.um.Unmarshal(operationResponse.Bytes())
	if unmarshalled != nil {
		contained = unmarshalled.(*stu3pb.ContainedResource)
	}
	return contained, resp, doErr
}
