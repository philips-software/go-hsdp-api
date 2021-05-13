package internal

import "encoding/json"

// Bundle represents a FHIR bundle response
type Bundle struct {
	Type  string        `json:"type,omitempty"`
	Total int64         `json:"total,omitempty"`
	Entry []BundleEntry `json:"entry,omitempty"`
}

// BundleEntry represents a entry item in a bundle
type BundleEntry struct {
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource json.RawMessage `json:"resource,omitempty"`
}
