package internal

import "encoding/json"

// Bundle represents a FHIR bundle response
type Bundle struct {
	Type  string        `json:"type,omitempty"`
	Total int64         `json:"total,omitempty"`
	Entry []BundleEntry `json:"entry,omitempty"`
	Link  BundleLinks   `json:"link,omitempty"`
}

type BundleLinks []LinkURL

func (b *BundleLinks) Next() *LinkURL {
	if b == nil {
		return nil
	}
	for _, e := range *b {
		if e.Relation == "next" {
			return &e
		}
	}
	return nil
}

type LinkURL struct {
	URL      string `json:"url"`
	Relation string `json:"relation"`
}

// BundleEntry represents a entry item in a bundle
type BundleEntry struct {
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource json.RawMessage `json:"resource,omitempty"`
}
