package identifier

import (
	"strings"

	r4gp "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
)

func UseToString(val *r4dt.Identifier_UseCode) string {
	enum := val.Value.Enum()
	if enum != nil {
		return enum.String()
	}
	return ""
}

func StringToUse(use string) *r4dt.Identifier_UseCode {
	switch strings.ToLower(use) {
	case "temp":
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_TEMP,
		}
	case "usual":
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_USUAL,
		}
	case "official":
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_OFFICIAL,
		}
	case "old":
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_OLD,
		}
	case "secondary":
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_SECONDARY,
		}
	default:
		return &r4dt.Identifier_UseCode{
			Value: r4gp.IdentifierUseCode_INVALID_UNINITIALIZED,
		}
	}
}
