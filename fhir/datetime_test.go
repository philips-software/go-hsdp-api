package fhir

import (
	"encoding/json"
	"testing"
)

type Simple struct {
	Foo []FHIRDateTime `bson:"foo,omitempty" json:"foo,omitempty"`
}

func TestFHIRDateTime(t *testing.T) {
	simple := &Simple{}

	data := []byte("{ \"foo\": [\"1991-02-01T10:00:00-05:00\", \"1992-02-01\", \"1993-02-01T10:00:00-05:00\"]}")
	err := json.Unmarshal(data, &simple)
	if err != nil {
		t.Errorf("Error unmarhsaling: %v", err)
		return
	}
	if len(simple.Foo) != 3 {
		t.Errorf("Expected 3 entries, got: %d", len(simple.Foo))
	}

	simple.Foo[0].Precision = Timestamp
	_, err = json.Marshal(&simple)
	if err != nil {
		t.Errorf("Error marshaling: %v", err)
		return
	}
}
