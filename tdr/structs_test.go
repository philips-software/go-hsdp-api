package tdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataType(t *testing.T) {
	var dataType = DataType{
		System: "Go",
		Code:   "Test",
	}
	assert.Equal(t, "tdr.DataType:System=Go,Code=Test", dataType.String())
}
