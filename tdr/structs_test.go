package tdr

import (
	"testing"
)

func TestDataType(t *testing.T) {
	var dataType = DataType{
		System: "Go",
		Code:   "Test",
	}
	if dataType.String() != "tdr.DataType:System=Go,Code=Test" {
		t.Errorf("Unexpected output: %s", dataType.String())
	}
}
