package stu3

import (
	"time"

	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/codes_go_proto"

	stu3dt "github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

const (
	ExtDeleteURL = "http://hsdp.com/cdr/Subscription/deletionUri"
)

type WithSubscriptionFunc func(sub *stu3pb.Subscription) error
type StringValue func(sub *stu3pb.Subscription) string

func WithCriteria(critera string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		sub.Criteria = &stu3dt.String{Value: critera}
		return nil
	}
}

func WithReason(reason string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		sub.Reason = &stu3dt.String{Value: reason}
		return nil
	}
}

// DeleteEndpointValue returns the URI if set, empty string otherwise
func DeleteEndpointValue() StringValue {
	return func(sub *stu3pb.Subscription) string {
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
func WithDeleteEndpoint(endpoint string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		if endpoint == "" {
			return nil
		}
		if sub.Channel == nil {
			sub.Channel = &stu3pb.Subscription_Channel{}
		}
		if sub.Channel.Extension == nil {
			sub.Channel.Extension = make([]*stu3dt.Extension, 0)
		}
		sub.Channel.Extension = append(sub.Channel.Extension, &stu3dt.Extension{
			Url: &stu3dt.Uri{Value: ExtDeleteURL},
			Value: &stu3dt.Extension_ValueX{
				Choice: &stu3dt.Extension_ValueX_Uri{
					Uri: &stu3dt.Uri{Value: endpoint},
				},
			},
		})
		sub.Channel.Type = &codes_go_proto.SubscriptionChannelTypeCode{
			Value: codes_go_proto.SubscriptionChannelTypeCode_REST_HOOK,
		}
		sub.Channel.Payload = &stu3dt.String{Value: "application/fhir+json"}
		return nil
	}
}

func WithEndpoint(endpoint string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		if sub.Channel == nil {
			sub.Channel = &stu3pb.Subscription_Channel{}
		}
		sub.Channel.Endpoint = &stu3dt.Uri{
			Value: endpoint,
		}
		sub.Channel.Type = &codes_go_proto.SubscriptionChannelTypeCode{
			Value: codes_go_proto.SubscriptionChannelTypeCode_REST_HOOK,
		}
		sub.Channel.Payload = &stu3dt.String{Value: "application/fhir+json"}
		return nil
	}
}

func WithHeaders(headers []string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		if len(headers) == 0 {
			return nil
		}
		if sub.Channel == nil {
			sub.Channel = &stu3pb.Subscription_Channel{}
		}
		sub.Channel.Header = make([]*stu3dt.String, len(headers))
		for i, h := range headers {
			sub.Channel.Header[i] = &stu3dt.String{Value: h}
		}
		return nil
	}
}

func WithContact(system, value, use string) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		if sub.Contact == nil {
			sub.Contact = make([]*stu3dt.ContactPoint, 0)
		}
		rank := len(sub.Contact) + 1
		useCode := stu3dt.ContactPointUseCode_INVALID_UNINITIALIZED
		useSystem := stu3dt.ContactPointSystemCode_INVALID_UNINITIALIZED

		switch use {
		case "work":
			useCode = stu3dt.ContactPointUseCode_WORK
		case "home":
			useCode = stu3dt.ContactPointUseCode_HOME
		case "mobile":
			useCode = stu3dt.ContactPointUseCode_MOBILE
		case "old":
			useCode = stu3dt.ContactPointUseCode_OLD
		case "temp":
			useCode = stu3dt.ContactPointUseCode_TEMP
		}
		switch system {
		case "email":
			useSystem = stu3dt.ContactPointSystemCode_EMAIL
		case "phone":
			useSystem = stu3dt.ContactPointSystemCode_PHONE
		case "fax":
			useSystem = stu3dt.ContactPointSystemCode_FAX
		case "pager":
			useSystem = stu3dt.ContactPointSystemCode_PAGER
		case "url":
			useSystem = stu3dt.ContactPointSystemCode_URL
		case "sms":
			useSystem = stu3dt.ContactPointSystemCode_SMS
		case "other":
			useSystem = stu3dt.ContactPointSystemCode_OTHER
		}
		sub.Contact = append(sub.Contact, &stu3dt.ContactPoint{
			Rank:   &stu3dt.PositiveInt{Value: uint32(rank)},
			Use:    &stu3dt.ContactPointUseCode{Value: useCode},
			Value:  &stu3dt.String{Value: value},
			System: &stu3dt.ContactPointSystemCode{Value: useSystem},
		})
		return nil
	}
}

// WithEndtime sets the end time of the subscription
func WithEndtime(at time.Time) WithSubscriptionFunc {
	return func(sub *stu3pb.Subscription) error {
		sub.End = &stu3dt.Instant{
			Precision: stu3dt.Instant_MICROSECOND,
			ValueUs:   at.UnixNano() / 1000,
		}
		return nil
	}
}

// NewSubscription creates a FHIR Subscription proto resource
// The WithSubscriptionFunc option methods should be used to build the structure
func NewSubscription(options ...WithSubscriptionFunc) (*stu3pb.Subscription, error) {
	sub := &stu3pb.Subscription{}
	sub.Status = &codes_go_proto.SubscriptionStatusCode{
		Value: codes_go_proto.SubscriptionStatusCode_REQUESTED,
	}
	for _, w := range options {
		if err := w(sub); err != nil {
			return nil, err
		}
	}
	return sub, nil
}
