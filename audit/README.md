## Using the audit client

```go
package main

import (
        "net/http"
        "fmt"
        "time"
        "github.com/philips-software/go-hsdp-api/audit"

	dstu2cp "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/codes_go_proto"
	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"
	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"
)

func main() {
        client, err := audit.NewClient(http.DefaultClient, logging.Config{
                SharedKey:    "YourSharedKeyHere=",
                SharedSecret: "YourSharedSecretHere==",
                BaseURL:      "https://audit-service.host.com",
        })
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
	timestamp := time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)
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
        _, _, err = client.CreateAuditEvent(event)
        if err != nil {
            fmt.Printf("Create audit event failed: %v\n", err)
        }
}
```
