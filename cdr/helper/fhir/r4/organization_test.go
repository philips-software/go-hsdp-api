package r4_test

import (
	"testing"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4"
	"github.com/stretchr/testify/assert"
)

func TestNewOrganization(t *testing.T) {
	org, err := r4.NewOrganization("Europe/Amsterdam", "id-here", "Hospital")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, org) {
		return
	}
	assert.Equal(t, "Hospital", org.Name.Value)
	assert.Equal(t, "id-here", org.Identifier[0].GetValue().Value)
}
