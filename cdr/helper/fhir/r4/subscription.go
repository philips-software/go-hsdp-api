package r4

import (
	"time"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"

	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pbsub "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/subscription_go_proto"
)

const (
	ExtDeleteURL = "http://hsdp.com/cdr/4.0/Subscription/deletionUri"
)

type WithFunc func(sub *r4pbsub.Subscription) error
type StringValue func(sub *r4pbsub.Subscription) string

func WithCriteria(critera string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		sub.Criteria = &r4dt.String{Value: critera}
		return nil
	}
}

func WithReason(reason string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		sub.Reason = &r4dt.String{Value: reason}
		return nil
	}
}

// DeleteEndpointValue returns the URI if set, empty string otherwise
func DeleteEndpointValue() StringValue {
	return func(sub *r4pbsub.Subscription) string {
		if sub.Channel == nil {
			return ""
		}
		if sub.Channel.Extension == nil {
			return ""
		}
		if len(sub.Channel.Extension) == 0 {
			return ""
		}
		for _, e := range sub.Channel.Extension {
			if e.Url.Value != ExtDeleteURL || e.Value == nil {
				continue
			}
			uri := e.Value.GetUri()
			if uri == nil {
				continue
			}
			return uri.Value

		}
		return ""
	}
}

// WithDeleteEndpoint adds an endpoint which is called a Resource is deleted
// This is an extension supported by CDR
func WithDeleteEndpoint(endpoint string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		if endpoint == "" {
			return nil
		}
		if sub.Channel == nil {
			sub.Channel = &r4pbsub.Subscription_Channel{}
		}
		if sub.Channel.Extension == nil {
			sub.Channel.Extension = make([]*r4dt.Extension, 0)
		}
		sub.Channel.Extension = append(sub.Channel.Extension, &r4dt.Extension{
			Url: &r4dt.Uri{Value: ExtDeleteURL},
			Value: &r4dt.Extension_ValueX{
				Choice: &r4dt.Extension_ValueX_Uri{
					Uri: &r4dt.Uri{Value: endpoint},
				},
			},
		})
		sub.Channel.Type = &r4pbsub.Subscription_Channel_TypeCode{
			Value: codes_go_proto.SubscriptionChannelTypeCode_REST_HOOK,
		}
		sub.Channel.Payload = &r4pbsub.Subscription_Channel_PayloadCode{
			Value: "application/fhir+json",
		}

		return nil
	}
}

func WithEndpoint(endpoint string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		if sub.Channel == nil {
			sub.Channel = &r4pbsub.Subscription_Channel{}
		}
		sub.Channel.Endpoint = &r4dt.Url{
			Value: endpoint,
		}
		sub.Channel.Type = &r4pbsub.Subscription_Channel_TypeCode{
			Value: codes_go_proto.SubscriptionChannelTypeCode_REST_HOOK,
		}
		sub.Channel.Payload = &r4pbsub.Subscription_Channel_PayloadCode{
			Value: "application/fhir+json;fhirVersion=4.0",
		}
		return nil
	}
}

func WithHeaders(headers []string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		if len(headers) == 0 {
			return nil
		}
		if sub.Channel == nil {
			sub.Channel = &r4pbsub.Subscription_Channel{}
		}
		sub.Channel.Header = make([]*r4dt.String, len(headers))
		for i, h := range headers {
			sub.Channel.Header[i] = &r4dt.String{Value: h}
		}
		return nil
	}
}

func WithContact(system, value, use string) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		if sub.Contact == nil {
			sub.Contact = make([]*r4dt.ContactPoint, 0)
		}
		rank := len(sub.Contact) + 1
		useCode := codes_go_proto.ContactPointUseCode_INVALID_UNINITIALIZED
		useSystem := codes_go_proto.ContactPointSystemCode_INVALID_UNINITIALIZED

		switch use {
		case "work":
			useCode = codes_go_proto.ContactPointUseCode_WORK
		case "home":
			useCode = codes_go_proto.ContactPointUseCode_HOME
		case "mobile":
			useCode = codes_go_proto.ContactPointUseCode_MOBILE
		case "old":
			useCode = codes_go_proto.ContactPointUseCode_OLD
		case "temp":
			useCode = codes_go_proto.ContactPointUseCode_TEMP
		}
		switch system {
		case "email":
			useSystem = codes_go_proto.ContactPointSystemCode_EMAIL
		case "phone":
			useSystem = codes_go_proto.ContactPointSystemCode_PHONE
		case "fax":
			useSystem = codes_go_proto.ContactPointSystemCode_FAX
		case "pager":
			useSystem = codes_go_proto.ContactPointSystemCode_PAGER
		case "url":
			useSystem = codes_go_proto.ContactPointSystemCode_URL
		case "sms":
			useSystem = codes_go_proto.ContactPointSystemCode_SMS
		case "other":
			useSystem = codes_go_proto.ContactPointSystemCode_OTHER
		}
		sub.Contact = append(sub.Contact, &r4dt.ContactPoint{
			Rank:   &r4dt.PositiveInt{Value: uint32(rank)},
			Use:    &r4dt.ContactPoint_UseCode{Value: useCode},
			Value:  &r4dt.String{Value: value},
			System: &r4dt.ContactPoint_SystemCode{Value: useSystem},
		})
		return nil
	}
}

// WithEndtime sets the end time of the subscription
func WithEndtime(at time.Time) WithFunc {
	return func(sub *r4pbsub.Subscription) error {
		sub.End = &r4dt.Instant{
			Precision: r4dt.Instant_MICROSECOND,
			ValueUs:   at.UnixNano() / 1000,
		}
		return nil
	}
}

// NewSubscription creates a FHIR Subscription proto resource
// The WithFunc option methods should be used to build the structure
func NewSubscription(options ...WithFunc) (*r4pbsub.Subscription, error) {
	sub := &r4pbsub.Subscription{}
	sub.Status = &r4pbsub.Subscription_StatusCode{
		Value: codes_go_proto.SubscriptionStatusCode_REQUESTED,
	}
	for _, w := range options {
		if err := w(sub); err != nil {
			return nil, err
		}
	}
	return sub, nil
}
