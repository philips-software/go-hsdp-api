package iam

// Permission represents a IAM Permission resource
type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Type        string `json:"type"`
}
