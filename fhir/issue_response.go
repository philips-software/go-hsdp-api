package fhir

// Code describes a coding
type Code struct {
	Coding []Coding `json:"coding"`
}

// Issue descrbies an issue
type Issue struct {
	Severity string `json:"Severity"`
	Details  string `json:"Details"`
	Code     Code   `json:"Code"`
}

// IssueResponse enscapsaltes one or more issues
type IssueResponse struct {
	Issues []Issue `json:"issue"`
}
