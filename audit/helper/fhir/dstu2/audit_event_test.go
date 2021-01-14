package dstu2_test

import (
	"testing"

	"github.com/philips-software/go-hsdp-api/audit/helper/fhir/dstu2"
	"github.com/stretchr/testify/assert"
)

func TestNewAuditEvent(t *testing.T) {
	event, err := dstu2.NewAuditEvent("key", "tenant")

	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, event) {
		return
	}
}
