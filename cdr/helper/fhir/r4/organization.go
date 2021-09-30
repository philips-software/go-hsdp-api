// Package r4 contains helper methods for use with CDR
package r4

import (
	"encoding/json"
	"fmt"

	"github.com/google/fhir/go/jsonformat"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	r4pborg "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
)

// NewOrganization returns a CDR R4 organization in Google FHIR proto format
func NewOrganization(timeZone, orgID, name string) (*r4pborg.Organization, error) {
	org := map[string]interface{}{
		"resourceType": "Organization",
		"id":           orgID,
		"name":         name,
		"identifier": []map[string]interface{}{
			{
				"use":    "usual",
				"system": "https://identity.philips-healthsuite.com/organization",
				"value":  orgID,
			},
		},
	}
	jsonPayload, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	um, err := jsonformat.NewUnmarshaller(timeZone, jsonformat.R4)
	if err != nil {
		return nil, fmt.Errorf("failed to create unmarshaller %v", err)
	}
	unmarshalled, err := um.Unmarshal(jsonPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal organization: %v", err)
	}
	contained := unmarshalled.(*r4pb.ContainedResource)
	organization := contained.GetOrganization()
	return organization, nil
}
