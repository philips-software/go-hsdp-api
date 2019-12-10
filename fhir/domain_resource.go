package fhir

type DomainResource struct {
	Resource
	Text      *Narrative    `json:"text,omitempty"`
	Contained []interface{} `json:"contained,omitempty"`
}
