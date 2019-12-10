package fhir

import (
	"testing"
)

func TestParseError(t *testing.T) {
	var err1 = []interface{}{
		"Foo",
		"Bar",
	}
	var err3 = map[string]interface{}{
		"Foo": "Bar",
	}
	parsed := ParseError(err1)
	if parsed != `[Foo, Bar]` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}
	parsed = ParseError("err2")
	if parsed != `err2` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}
	parsed = ParseError(err3)
	if parsed != `{Foo: Bar}` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}

}
