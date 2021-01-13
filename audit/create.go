package audit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"
)

func (c *Client) CreateAuditEvent(event *dstu2pb.AuditEvent) (*dstu2pb.ContainedResource, *Response, error) {
	eventJSON, err := c.ma.MarshalResource(event)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.NewAuditRequest("POST", "core/audit/AuditEvent", eventJSON, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("audit.CreateAuditEvent: %w", err)
	}
	if err := c.httpSigner.SignRequest(req); err != nil {
		return nil, nil, err
	}
	var operationResponse bytes.Buffer
	resp, err := c.Do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	// Success
	if resp.StatusCode == http.StatusCreated {
		return &dstu2pb.ContainedResource{}, resp, nil
	}
	// OperationOutcome
	unmarshalled, err := c.um.Unmarshal(operationResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*dstu2pb.ContainedResource)
	return contained, resp, nil
}
