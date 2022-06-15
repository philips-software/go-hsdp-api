package identifier_test

import (
	"testing"

	r4gp "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4/identifier"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierToString(t *testing.T) {
	val := identifier.UseToString(&r4dt.Identifier_UseCode{
		Value: r4gp.IdentifierUseCode_TEMP,
	})
	assert.Equal(t, "TEMP", val)
}
