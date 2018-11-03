package tdr

import "encoding/json"

type DeletePolicy struct {
	Duration int    `json:"duration,omitempty"`
	Unit     string `json:"unit,omitempty"`
}

// Contract describes a TDR Contract
type Contract struct {
	ID                         string          `json:"id,omitempty"`
	Meta                       *Meta           `json:"meta,omitempty"`
	DataType                   DataType        `json:"dataType,omitempty"`
	Schema                     json.RawMessage `json:"schema,omitempty"`
	Organization               string          `json:"organization"`
	SendNotifications          bool            `json:"sendNotifications"`
	NotificationServiceTopicID string          `json:"notificationServiceTopicId,omitempty"`
	DeletePolicy               DeletePolicy    `json:"deletePolicy"`
}
