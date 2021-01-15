## Using the audit client

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	dstu2ct "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/codes_go_proto"
	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"
	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"

	"github.com/philips-software/go-hsdp-api/audit/helper/fhir/dstu2"

	"github.com/philips-software/go-hsdp-api/audit"
)

func main() {
	productKey := "xxx-your-key-here-xxx"
	now := time.Now()

	client, err := audit.NewClient(http.DefaultClient, &audit.Config{
		SharedSecret: "secrethere",
		SharedKey:    "keyhere",
		AuditBaseURL: "https://your-create-url-here.eu-west.philips-healthsuite.com",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	event, err := dstu2.NewAuditEvent(productKey, "andy",
		dstu2.AddSourceExtensionUriValue("applicationName", "patientapp"),
		dstu2.AddParticipant(&dstu2pb.AuditEvent_Participant{
			UserId: &dstu2dt.Identifier{
				Value: &dstu2dt.String{Value: "smokeuser@philips.com"},
			},
			Requestor: &dstu2dt.Boolean{Value: true},
		}),
		dstu2.WithEvent(&dstu2pb.AuditEvent_Event{
			Action: &dstu2ct.AuditEventActionCode{
				Value: dstu2ct.AuditEventActionCode_E,
			},
			DateTime: dstu2.DateTime(now),
			Type: &dstu2dt.Coding{
				System:  &dstu2dt.Uri{Value: "http://hl7.org/fhir/ValueSet/audit-event-type"},
				Version: &dstu2dt.String{Value: "1"},
				Code:    &dstu2dt.Code{Value: "11011"},
				Display: &dstu2dt.String{Value: fmt.Sprintf("Timestamp %v", now.String())},
			},
			Outcome: &dstu2ct.AuditEventOutcomeCode{
				Value: dstu2ct.AuditEventOutcomeCode_INVALID_UNINITIALIZED,
			},
			OutcomeDesc: &dstu2dt.String{Value: "Success"},
		}),
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
		}))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	outcome, resp, err := client.CreateAuditEvent(event)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if resp == nil {
		fmt.Printf("response is nil\n")
		return
	}
	fmt.Printf("Audit result: %d\n", resp.StatusCode)
	if outcome != nil {
		fmt.Printf("Outcome: %v\n", outcome)
	}
}
```
