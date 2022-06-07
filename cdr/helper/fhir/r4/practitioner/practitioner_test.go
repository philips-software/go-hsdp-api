package practitioner_test

import (
	"testing"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4/practitioner"
	"github.com/stretchr/testify/assert"
)

func TestNewPractitioner(t *testing.T) {
	p, err := practitioner.NewPractitioner(
		practitioner.WithIdentifier(
			"https://iam-client-test.us-east.philips-healthsuite.com/oauth2/access_token",
			"ron.swanson@pawnee.io",
			"temp"),
		practitioner.WithName("Ron Swanson", "Swanson", []string{"Ron"}),
	)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, p) {
		return
	}
	if !assert.Len(t, p.Identifier, 1) {
		return
	}
	assert.Equal(t, "ron.swanson@pawnee.io", p.Identifier[0].Value.GetValue())
	assert.Equal(t, "Swanson", p.Name[0].Family.GetValue())
}
