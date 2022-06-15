package identifier_test

import (
	"testing"

	stu3dt "github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3/identifier"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierToString(t *testing.T) {
	val := identifier.UseToString(&stu3dt.IdentifierUseCode{
		Value: stu3dt.IdentifierUseCode_TEMP,
	})
	assert.Equal(t, "TEMP", val)
}
