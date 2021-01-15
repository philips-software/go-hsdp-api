package dstu2_test

import (
	"testing"
	"time"

	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"

	dstu2ct "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/codes_go_proto"
	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"

	"github.com/philips-software/go-hsdp-api/audit/helper/fhir/dstu2"
	"github.com/stretchr/testify/assert"
)

func TestNewAuditEvent(t *testing.T) {
	event, err := dstu2.NewAuditEvent("key", "tenant",
		dstu2.WithSourceIdentifier(&dstu2dt.Identifier{
			Value: &dstu2dt.String{Value: "smokeuser@philips.com"},
			Type: &dstu2dt.CodeableConcept{
				Coding: []*dstu2dt.Coding{
					{
						System:  &dstu2dt.Uri{Value: "http://hl7.org/fhir/ValueSet/identifier-type"},
						Code:    &dstu2dt.Code{Value: "4"},
						Display: &dstu2dt.String{Value: "application server"},
					},
				},
			},
		}),
		dstu2.WithEvent(&dstu2pb.AuditEvent_Event{
			Action: &dstu2ct.AuditEventActionCode{
				Value: dstu2ct.AuditEventActionCode_E,
			},
			DateTime: dstu2.DateTime(time.Now()),
			Type: &dstu2dt.Coding{
				System:  &dstu2dt.Uri{Value: "http://hl7.org/fhir/ValueSet/audit-event-type"},
				Version: &dstu2dt.String{Value: "1"},
				Code:    &dstu2dt.Code{Value: "11011"},
				Display: &dstu2dt.String{Value: "Testing"},
			},
			Outcome: &dstu2ct.AuditEventOutcomeCode{
				Value: dstu2ct.AuditEventOutcomeCode_INVALID_UNINITIALIZED,
			},
			OutcomeDesc: &dstu2dt.String{Value: "Success"},
		}),
		dstu2.AddParticipant(&dstu2pb.AuditEvent_Participant{
			UserId: &dstu2dt.Identifier{
				Value: &dstu2dt.String{Value: "smokeuser@philips.com"},
			},
			Requestor: &dstu2dt.Boolean{Value: true},
		}))

	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, event) {
		return
	}
}
