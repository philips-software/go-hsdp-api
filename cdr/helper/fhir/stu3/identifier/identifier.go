package identifier

import (
	stu3dt "github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
)

func UseToString(val *stu3dt.IdentifierUseCode) string {
	enum := val.Value.Enum()
	if enum != nil {
		return enum.String()
	}
	return ""
}

func StringToUse(use string) *stu3dt.IdentifierUseCode {
	switch use {
	case "temp":
		return &stu3dt.IdentifierUseCode{
			Value: stu3dt.IdentifierUseCode_TEMP,
		}
	case "usual":
		return &stu3dt.IdentifierUseCode{
			Value: stu3dt.IdentifierUseCode_USUAL,
		}
	case "official":
		return &stu3dt.IdentifierUseCode{
			Value: stu3dt.IdentifierUseCode_OFFICIAL,
		}
	case "secondary":
		return &stu3dt.IdentifierUseCode{
			Value: stu3dt.IdentifierUseCode_SECONDARY,
		}
	default:
		return &stu3dt.IdentifierUseCode{
			Value: stu3dt.IdentifierUseCode_INVALID_UNINITIALIZED,
		}
	}
}
