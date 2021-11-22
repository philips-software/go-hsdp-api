package internal

type OperationOutcome struct {
	Issue        []Issue `json:"issue"`
	ResourceType string  `json:"resourceType"`
}

type Issue struct {
	Severity    string  `json:"severity"`
	Code        string  `json:"code"`
	Details     Details `json:"details"`
	Diagnostics string  `json:"diagnostics"`
}

type Details struct {
	Coding Coding `json:"coding"`
	Text   string `json:"text"`
}

type Coding struct {
	System string `json:"system"`
	Code   string `json:"code"`
}
