package fhir

type Narrative struct {
	Status string `bson:"status,omitempty" json:"status,omitempty"`
	Div    string `bson:"div,omitempty" json:"div,omitempty"`
}
