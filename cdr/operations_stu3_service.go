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

// Patch makes changes to a FHIR resources accepting the JSONPatch format set
func (o *OperationsSTU3Service) Patch(resourceID string, jsonPatch []byte, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodPatch, resourceID, jsonPatch, options)
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

// Post creates new FHIR resources
func (o *OperationsSTU3Service) Post(resourceID string, jsonBody []byte, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	return o.postOrPut(http.MethodPost, resourceID, jsonBody, options...)
}

// Put creates or updates new FHIR resources
func (o *OperationsSTU3Service) Put(resourceID string, jsonBody []byte, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	return o.postOrPut(http.MethodPut, resourceID, jsonBody, options...)
}

// Get returns a FHIR resource
func (o *OperationsSTU3Service) Get(resourceID string, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodGet, resourceID, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/fhir+json")
	var operationResponse bytes.Buffer
	resp, err := o.client.Do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	unmarshalled, err := o.um.Unmarshal(operationResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	return contained, resp, nil
}

// Delete removes a FHIR resource
func (o *OperationsSTU3Service) Delete(resourceID string, options ...OptionFunc) (bool, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodDelete, resourceID, nil, options)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/fhir+json")
	var operationResponse bytes.Buffer
	resp, err := o.client.Do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}

func (o *OperationsSTU3Service) postOrPut(method, resourceID string, jsonBody []byte, options ...OptionFunc) (*stu3pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(method, resourceID, jsonBody, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/fhir+json")
	var operationResponse bytes.Buffer
	resp, err := o.client.Do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = ErrEmptyResult
		}
		return nil, resp, err
	}
	unmarshalled, err := o.um.Unmarshal(operationResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	return contained, resp, nil
}
