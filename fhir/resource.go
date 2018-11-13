package fhir

type Resource struct {
	ResourceType  string `json:"resourceType,omitempty"`
	ID            string `json:"id,omitempty"`
	Meta          *Meta  `json:"meta,omitempty"`
	ImplicitRules string `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	Language      string `bson:"language,omitempty" json:"language,omitempty"`
}
