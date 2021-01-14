package dstu2

import (
	dstu2dt "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/datatypes_go_proto"
	dstu2pb "github.com/google/fhir/go/proto/google/fhir/proto/dstu2/resources_go_proto"
)

type WithFunc func(sub *dstu2pb.AuditEvent) error

// NewAuditEvent creates a new audit event. It takes
// productKey and tenant as arguments as these are required
// for publishing to to Host Auditing service
func NewAuditEvent(productKey, tenant string, options ...WithFunc) (*dstu2pb.AuditEvent, error) {
	event := &dstu2pb.AuditEvent{}

	if err := WithSourceExtensionUriValue("productKey", productKey)(event); err != nil {
		return nil, err
	}
	if err := WithSourceExtensionUriValue("tenant", tenant)(event); err != nil {
		return nil, err
	}
	for _, w := range options {
		if err := w(event); err != nil {
			return nil, err
		}
	}
	return event, nil
}

// WithObject adds the passed object to the AuditEvent
func WithObject(object *dstu2pb.AuditEvent_Object) WithFunc {
	return func(event *dstu2pb.AuditEvent) error {
		if event.Object == nil {
			event.Object = []*dstu2pb.AuditEvent_Object{}
		}
		event.Object = append(event.Object, object)
		return nil
	}
}

// WithExtensionUriValue sets extension Uri/Value tuples, some of which are mandatory
// for successfully posting to HSDP Audit
func WithSourceExtensionUriValue(extensionUri, extensionValue string) WithFunc {
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
				Url: &dstu2dt.Uri{},
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
