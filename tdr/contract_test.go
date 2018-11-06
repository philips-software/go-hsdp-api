package tdr

import (
	"testing"
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
	if contract.String() != "tdr.Contract:ID=Contract,DataType=tdr.DataType:System=Go,Code=Test,Organization=TDROrg" {
		t.Errorf("Unexpected output: %s", contract.String())
	}
}
