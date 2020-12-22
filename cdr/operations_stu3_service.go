package cdr

import (
	"bytes"
	"io"
	"net/http"

	"github.com/google/fhir/go/jsonformat"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

type OperationsSTU3Service struct {
	client   *Client
	timeZone string
	ma       *jsonformat.Marshaller
	um       *jsonformat.Unmarshaller
}

func (o *OperationsSTU3Service) Patch(resourceID string, jsonPatch []byte, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	req, err := o.client.NewCDRRequest(http.MethodPatch, resourceID, jsonPatch, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json-patch+json")
	var patchResponse bytes.Buffer
	resp, err := o.client.Do(req, &patchResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	unmarshalled, err := o.um.Unmarshal(patchResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	return contained, resp, nil
}
