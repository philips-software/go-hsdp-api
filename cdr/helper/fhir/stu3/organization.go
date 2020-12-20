package stu3

import (
	"encoding/json"
	"fmt"

	"github.com/google/fhir/go/jsonformat"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

// NewOrganization returns a CDR STU3 organization in Google FHIR proto format
func NewOrganization(timeZone, orgID, name string) (*stu3pb.Organization, error) {
	org := map[string]interface{}{
		"resourceType": "Organization",
		"id":           orgID,
		"name":         name,
		"identifier": []map[string]interface{}{
			{
				"use":    "official",
				"system": "https://identity.philips-healthsuite.com/organization",
				"value":  orgID,
			},
		},
	}
	jsonPayload, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	um, err := jsonformat.NewUnmarshaller(timeZone, jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("failed to create unmarshaller %v", err)
	}
	unmarshalled, err := um.Unmarshal(jsonPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal organization: %v", err)
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	organization := contained.GetOrganization()
	return organization, nil
}
