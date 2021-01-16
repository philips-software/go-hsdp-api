package tpns

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
	parsed := parseError(err1)
	if parsed != `[Foo, Bar]` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}
	parsed = parseError("err2")
	if parsed != `err2` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}
	parsed = parseError(err3)
	if parsed != `{Foo: Bar}` {
		t.Errorf("Unexpected parse output: `%s`", parsed)
	}

}
