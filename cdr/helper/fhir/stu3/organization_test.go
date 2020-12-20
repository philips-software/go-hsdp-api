package stu3_test

import (
	"testing"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
	"github.com/stretchr/testify/assert"
)

func TestNewOrganization(t *testing.T) {
	org, err := stu3.NewOrganization("Europe/Amsterdam", "id-here", "Hospital")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, org) {
		return
	}
	assert.Equal(t, "Hospital", org.Name.Value)
	assert.Equal(t, "id-here", org.Identifier[0].GetValue().Value)
}
