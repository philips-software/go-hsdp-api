// Package stu3 contains helper methods for use with CDR
package stu3

import (
	"encoding/json"
	"fmt"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	rpb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

// NewOrganization returns a CDR STU3 organization in Google FHIR proto format
func NewOrganization(timeZone, orgID, name string) (*rpb.Organization, error) {
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

	um, err := jsonformat.NewUnmarshaller(timeZone, fhirversion.STU3)
	if err != nil {
		return nil, fmt.Errorf("failed to create unmarshaller %v", err)
	}

	contained, err := um.UnmarshalR3(jsonPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal organization: %v", err)
	}
	organization := contained.GetOrganization()
	return organization, nil
}
