package stu3

import (
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/codes_go_proto"

	stu3dt "github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

type WithFunc func(sub *stu3pb.Subscription) error

func WithCriteria(critera string) WithFunc {
	return func(sub *stu3pb.Subscription) error {
		sub.Criteria = &stu3dt.String{Value: critera}
		return nil
	}
}

func WithReason(reason string) WithFunc {
	return func(sub *stu3pb.Subscription) error {
		sub.Reason = &stu3dt.String{Value: reason}
		return nil
	}
}

func WithEndpoint(endpoint string) WithFunc {
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

func WithHeaders(headers []string) WithFunc {
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

func WithContact(system, value, use string) WithFunc {
	return func(sub *stu3pb.Subscription) error {
		if sub.Contact == nil {
			sub.Contact = make([]*stu3dt.ContactPoint, 0)
		}
		rank := len(sub.Contact) + 1
		useCode := stu3dt.ContactPointUseCode_HOME
		useSystem := stu3dt.ContactPointSystemCode_EMAIL

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
		default:
			useCode = stu3dt.ContactPointUseCode_INVALID_UNINITIALIZED
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
		default:
			useSystem = stu3dt.ContactPointSystemCode_INVALID_UNINITIALIZED
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

// NewSubscription creates a FHIR Subscription proto resource
// The WithFunc option methods should be used to build the structure
func NewSubscription(options ...WithFunc) (*stu3pb.Subscription, error) {
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
