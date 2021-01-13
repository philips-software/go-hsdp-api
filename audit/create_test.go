package audit_test

import (
	"net/http"
	"testing"
	"time"

	dstu2cp "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/codes_go_proto"
	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"
	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"

	"github.com/philips-software/go-hsdp-api/audit"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxAudit.HandleFunc("/core/audit/AuditEvent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			if !assert.Equal(t, audit.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	timestamp := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
	event := &dstu2pb.AuditEvent{
		Id: &dstu2dt.Id{Value: "someID"},
		Event: &dstu2pb.AuditEvent_Event{
			Type: &dstu2dt.Coding{
				System:  &dstu2dt.Uri{Value: "http://hl7.org/fhir/ValueSet/audit-event-type"},
				Version: &dstu2dt.String{Value: "1"},
				Code:    &dstu2dt.Code{Value: "110112"},
				Display: &dstu2dt.String{Value: "Query"},
			},
			Action: &dstu2cp.AuditEventActionCode{Value: dstu2cp.AuditEventActionCode_E},
			DateTime: &dstu2dt.Instant{
				Precision: dstu2dt.Instant_MICROSECOND,
				ValueUs:   timestamp.UnixNano() / 1000,
			},
			Outcome:     &dstu2cp.AuditEventOutcomeCode{Value: dstu2cp.AuditEventOutcomeCode_SUCCESS},
			OutcomeDesc: &dstu2dt.String{Value: "Success"},
		},
		Participant: []*dstu2pb.AuditEvent_Participant{
			{
				UserId:    &dstu2dt.Identifier{Value: &dstu2dt.String{Value: "smokeuser@philips.com"}},
				Requestor: &dstu2dt.Boolean{Value: true},
			},
		},
		Extension: []*dstu2dt.Extension{
			{
				Url: &dstu2dt.Uri{Value: "http://foo.bar/com"},
				Value: &dstu2dt.Extension_ValueX{
					Choice: &dstu2dt.Extension_ValueX_Uri{
						Uri: &dstu2dt.Uri{Value: "http://lala.local"},
					},
				},
			},
		},
	}
	contained, resp, err := auditClient.CreateAuditEvent(event)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, contained) {
		return
	}
}
