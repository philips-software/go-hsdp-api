package tdr

import "encoding/json"

// DataItem describes a TDR Data item
type DataItem struct {
	ID                string            `json:"id,omitempty"`
	Meta              Meta              `json:"meta,omitempty"`
	ResourceType      string            `json:"resourceType,omitempty"`
	Timestamp         string            `json:"timestamp,omitempty"`
	SequenceNumber    int               `json:"sequenceNumber,omitempty"`
	Device            Device            `json:"device,omitempty"`
	User              User              `json:"user,omitempty"`
	RelatedPeripheral RelatedPeripheral `json:"relatedPeripheral,omitempty"`
	RelatedUser       RelatedUser       `json:"relatedUser,omitempty"`
	DataType          DataType          `json:"dataType,omitempty"`
	Organization      string            `json:"organization,omitempty"`
	Application       string            `json:"application,omitempty"`
	Proposition       string            `json:"proposition,omitempty"`
	Subscription      string            `json:"subscription,omitempty"`
	DataSource        string            `json:"dataSource,omitempty"`
	DataCategory      string            `json:"dataCategory,omitempty"`
	Data              json.RawMessage   `json:"data,omitempty"`
	Blob              string            `json:"blob,omitempty"`
	DeleteTimestamp   string            `json:"deleteTimestamp,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
	Tombstone         bool              `json:"tombstone,omitempty"`
}

// Device describes a TDR device
type Device struct {
	Description string `json:"description,omitempty"`
	System      string `json:"system"`
	Value       string `json:"value"`
}

// User describes a TDR user
type User struct {
	Description string `json:"description,omitempty"`
	System      string `json:"system"`
	Value       string `json:"value"`
}

// RelatedPeripheral describes a TDR related peripheral
type RelatedPeripheral struct {
	Description string `json:"description,omitempty"`
	System      string `json:"system"`
	Value       string `json:"value"`
}

// RelatedUser describes a TDR related user
type RelatedUser struct {
	Description string `json:"description,omitempty"`
	System      string `json:"system"`
	Value       string `json:"value"`
}
