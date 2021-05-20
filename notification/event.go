package notification

import (
	"time"
)

// Event is an Amazon SNS HTTP message on which HSDP Notification is currently based on
type Event struct {
	Type             string    `json:"Type"`
	MessageID        string    `json:"MessageId"`
	Token            string    `json:"Token"`
	TopicARN         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	Subject          string    `json:"Subject,omitempty"`
	SubscribeURL     string    `json:"SubscribeURL,omitempty"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
}
