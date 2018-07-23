package tdr

import "encoding/json"

// Contract describes a TDR Contract
type Contract struct {
	ID   string `json:"id,omitempty"`
	Meta struct {
		LastUpdated string `json:"lastUpdated,omitempty"`
		VersionID   string `json:"versionId,omitempty"`
	} `json:"meta,omitempty"`
	DataType struct {
		System string `json:"system,omitempty"`
		Code   string `json:"code,omitempty"`
	} `json:"dataType,omitempty"`
	Schema                     json.RawMessage `json:"schema,omitempty"`
	Organization               string          `json:"organization"`
	SendNotifications          bool            `json:"sendNotifications"`
	NotificationServiceTopicID string          `json:"notificationServiceTopicId,omitempty"`
	DeletePolicy               struct {
		Duration int    `json:"duration,omitempty"`
		Unit     string `json:"unit,omitempty"`
	} `json:"deletePolicy,omitempty"`
}
