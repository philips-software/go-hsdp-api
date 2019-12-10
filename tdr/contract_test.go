package tdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContract(t *testing.T) {
	var contract = Contract{
		ID: "Contract",
		DataType: DataType{
			System: "Go",
			Code:   "Test",
		},
		Organization: "TDROrg",
	}
	assert.Equal(t, "tdr.Contract:ID=Contract,DataType=tdr.DataType:System=Go,Code=Test,Organization=TDROrg", contract.String())
}
