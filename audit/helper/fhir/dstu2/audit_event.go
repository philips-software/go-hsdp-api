// Package dstu2 contains helper methods to construct AuditEvent resources
package dstu2

import (
	"time"

	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"
	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"
)

type OptionFunc func(sub *dstu2pb.AuditEvent) error

// NewAuditEvent creates a new audit event. It takes
// productKey and tenant as arguments as these are required
// for publishing to to Host Auditing service
func NewAuditEvent(productKey, tenant string, options ...OptionFunc) (*dstu2pb.AuditEvent, error) {
	event := &dstu2pb.AuditEvent{}

	if err := AddSourceExtensionUriValue("productKey", productKey)(event); err != nil {
		return nil, err
	}
	if err := AddSourceExtensionUriValue("tenant", tenant)(event); err != nil {
		return nil, err
	}
	for _, w := range options {
		if err := w(event); err != nil {
			return nil, err
		}
	}
	return event, nil
}

// WithEvent sets the event
func WithEvent(e *dstu2pb.AuditEvent_Event) OptionFunc {
	return func(event *dstu2pb.AuditEvent) error {
		event.Event = e
		return nil
	}
}

// DateTime returns DateTime
func DateTime(at time.Time) *dstu2dt.Instant {
	return &dstu2dt.Instant{
		Precision: dstu2dt.Instant_MICROSECOND,
		ValueUs:   at.UnixNano() / 1000,
	}
}

// AddParticipant adds the participant
func AddParticipant(participant *dstu2pb.AuditEvent_Participant) OptionFunc {
	return func(event *dstu2pb.AuditEvent) error {
		if event.Participant == nil {
			event.Participant = []*dstu2pb.AuditEvent_Participant{}
		}
		event.Participant = append(event.Participant, participant)
		return nil
	}
}

// AddObject adds the passed object to the AuditEvent
func AddObject(object *dstu2pb.AuditEvent_Object) OptionFunc {
	return func(event *dstu2pb.AuditEvent) error {
		if event.Object == nil {
			event.Object = []*dstu2pb.AuditEvent_Object{}
		}
		event.Object = append(event.Object, object)
		return nil
	}
}

// WithSourceIdentifier sets the source identifier
func WithSourceIdentifier(identifier *dstu2dt.Identifier) OptionFunc {
	return func(event *dstu2pb.AuditEvent) error {
		if event.Source == nil {
			event.Source = &dstu2pb.AuditEvent_Source{}
		}
		event.Source.Identifier = identifier
		return nil
	}
}

// AddSourceExtensionUriValue sets extension Uri/Value tuples, some of which are mandatory
// for successfully posting to HSDP Audit
func AddSourceExtensionUriValue(extensionUri, extensionValue string) OptionFunc {
	return func(event *dstu2pb.AuditEvent) error {
		if event.Source == nil {
			event.Source = &dstu2pb.AuditEvent_Source{}
		}
		if event.Source.Extension == nil {
			event.Source.Extension = []*dstu2dt.Extension{}
		}
		var ext *dstu2dt.Extension
		// Find the extension
		for _, e := range event.Source.Extension {
			if e.Url != nil && e.Url.Value == "/fhir/device" {
				ext = e
				break
			}
		}
		if ext == nil {
			ext = &dstu2dt.Extension{
				Url: &dstu2dt.Uri{Value: "/fhir/device"},
			}
			event.Source.Extension = append(event.Source.Extension, ext)
		}
		if ext.Extension == nil {
			ext.Extension = []*dstu2dt.Extension{}
		}
		var extensionEntry *dstu2dt.Extension
		for _, e := range ext.Extension {
			if e.Url != nil && e.Url.Value == extensionUri {
				extensionEntry = e
				break
			}
		}
		if extensionEntry == nil {
			extensionEntry = &dstu2dt.Extension{
				Url: &dstu2dt.Uri{
					Value: extensionUri,
				},
			}
			ext.Extension = append(ext.Extension, extensionEntry)
		}
		extensionEntry.Value = &dstu2dt.Extension_ValueX{
			Choice: &dstu2dt.Extension_ValueX_StringValue{
				StringValue: &dstu2dt.String{Value: extensionValue},
			},
		}
		return nil
	}
}
