package practitioner

import (
	"fmt"

	stu3dt "github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
)

type WithFunc func(resource *stu3pb.Practitioner) error

func WithIdentifier(system, value, use string) WithFunc {
	return func(resource *stu3pb.Practitioner) error {
		if resource.Identifier == nil {
			resource.Identifier = make([]*stu3dt.Identifier, 0)
		}
		val := &stu3dt.Identifier{
			System: &stu3dt.Uri{Value: system},
			Value:  &stu3dt.String{Value: value},
		}
		switch use {
		case "":
			break
		case "temp":
			val.Use = &stu3dt.IdentifierUseCode{
				Value: stu3dt.IdentifierUseCode_TEMP,
			}
		case "usual":
			val.Use = &stu3dt.IdentifierUseCode{
				Value: stu3dt.IdentifierUseCode_USUAL,
			}
		case "official":
			val.Use = &stu3dt.IdentifierUseCode{
				Value: stu3dt.IdentifierUseCode_OFFICIAL,
			}
		case "secondary":
			val.Use = &stu3dt.IdentifierUseCode{
				Value: stu3dt.IdentifierUseCode_SECONDARY,
			}
		default:
			return fmt.Errorf("unknown use '%s'", use)
		}
		resource.Identifier = append(resource.Identifier, val)
		return nil
	}
}

func WithName(text, family string, given []string) WithFunc {
	return func(resource *stu3pb.Practitioner) error {
		if resource.Name == nil {
			resource.Name = make([]*stu3dt.HumanName, 0)
		}
		var givenList []*stu3dt.String
		for _, g := range given {
			givenList = append(givenList, &stu3dt.String{Value: g})
		}
		resource.Name = append(resource.Name, &stu3dt.HumanName{
			Text:   &stu3dt.String{Value: text},
			Given:  givenList,
			Family: &stu3dt.String{Value: family},
		})
		return nil
	}
}

// NewPractitioner creates a FHIR Practitioner proto resource
// The WithPractitionerFunc option methods should be used to build the structure
func NewPractitioner(options ...WithFunc) (*stu3pb.Practitioner, error) {
	resource := &stu3pb.Practitioner{}

	for _, w := range options {
		if err := w(resource); err != nil {
			return nil, err
		}
	}
	return resource, nil
}
