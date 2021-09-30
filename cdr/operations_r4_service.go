package cdr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/google/fhir/go/jsonformat"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
)

type OperationsR4Service struct {
	client   *Client
	timeZone string
	ma       *jsonformat.Marshaller
	um       *jsonformat.Unmarshaller
}

// Patch makes changes to a FHIR resources accepting the JSONPatch format set
func (o *OperationsR4Service) Patch(resourceID string, jsonPatch []byte, options ...OptionFunc) (*r4pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodPatch, resourceID, jsonPatch, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/json-patch+json")
	var patchResponse bytes.Buffer
	resp, err := o.client.do(req, &patchResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("OperationsSTU3Service.Patch: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	unmarshalled, err := o.um.Unmarshal(patchResponse.Bytes())
	if err != nil {
		return nil, resp, fmt.Errorf("FHIR unmarshal: %w", err)
	}
	contained := unmarshalled.(*r4pb.ContainedResource)
	return contained, resp, nil
}

// Post creates new FHIR resources
func (o *OperationsR4Service) Post(resourceID string, jsonBody []byte, options ...OptionFunc) (*r4pb.ContainedResource, *Response, error) {
	return o.postOrPut(http.MethodPost, resourceID, jsonBody, options...)
}

// Put creates or updates new FHIR resources
func (o *OperationsR4Service) Put(resourceID string, jsonBody []byte, options ...OptionFunc) (*r4pb.ContainedResource, *Response, error) {
	return o.postOrPut(http.MethodPut, resourceID, jsonBody, options...)
}

// Get returns a FHIR resource
func (o *OperationsR4Service) Get(resourceID string, options ...OptionFunc) (*r4pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodGet, resourceID, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
	var operationResponse bytes.Buffer
	resp, err := o.client.do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("OperationsSTU3Service.Get: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	unmarshalled, err := o.um.Unmarshal(operationResponse.Bytes())
	if err != nil {
		return nil, resp, fmt.Errorf("FHIR unmarshal: %w", err)
	}
	contained := unmarshalled.(*r4pb.ContainedResource)
	return contained, resp, nil
}

// Delete removes a FHIR resource
func (o *OperationsR4Service) Delete(resourceID string, options ...OptionFunc) (bool, *Response, error) {
	req, err := o.client.newCDRRequest(http.MethodDelete, resourceID, nil, options)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
	var operationResponse bytes.Buffer
	resp, err := o.client.do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("OperationsSTU3Service.Delete: %w", ErrEmptyResult)
		}
		return false, resp, err
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}

func (o *OperationsR4Service) postOrPut(method, resourceID string, jsonBody []byte, options ...OptionFunc) (*r4pb.ContainedResource, *Response, error) {
	req, err := o.client.newCDRRequest(method, resourceID, jsonBody, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/fhir+json;fhirVersion=4.0")
	req.Header.Set("Content-Type", "application/fhir+json;fhirVersion=4.0")
	var operationResponse bytes.Buffer
	resp, err := o.client.do(req, &operationResponse)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("OperationsR4Service %s: %w", method, ErrEmptyResult)
		}
		return nil, resp, err
	}
	if operationResponse.Len() == 0 { // Empty body
		return &r4pb.ContainedResource{}, resp, nil
	}
	unmarshalled, err := o.um.Unmarshal(operationResponse.Bytes())
	if err != nil {
		return nil, resp, fmt.Errorf("FHIR unmarshal: %w", err)
	}
	contained := unmarshalled.(*r4pb.ContainedResource)
	return contained, resp, nil
}
