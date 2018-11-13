package fhir

type OperationOutcome struct {
	DomainResource
	Issue []OperationOutcomeIssueComponent `json:"issue,omitempty"`
}

type OperationOutcomeIssueComponent struct {
	BackboneElement
	Severity    string           `json:"severity,omitempty"`
	Code        string           `json:"code,omitempty"`
	Details     *CodeableConcept `json:"details,omitempty"`
	Diagnostics string           `json:"diagnostics,omitempty"`
	Location    []string         `json:"location,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}
