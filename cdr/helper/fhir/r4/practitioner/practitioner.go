package practitioner

import (
	"fmt"

	r4gp "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pprac "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_go_proto"
)

type WithFunc func(resource *r4pprac.Practitioner) error

func WithIdentifier(system, value, use string) WithFunc {
	return func(resource *r4pprac.Practitioner) error {
		if resource.Identifier == nil {
			resource.Identifier = make([]*r4dt.Identifier, 0)
		}
		val := &r4dt.Identifier{
			System: &r4dt.Uri{Value: system},
			Value:  &r4dt.String{Value: value},
		}
		switch use {
		case "":
			break
		case "temp":
			val.Use = &r4dt.Identifier_UseCode{
				Value: r4gp.IdentifierUseCode_TEMP,
			}
		case "usual":
			val.Use = &r4dt.Identifier_UseCode{
				Value: r4gp.IdentifierUseCode_USUAL,
			}
		case "official":
			val.Use = &r4dt.Identifier_UseCode{
				Value: r4gp.IdentifierUseCode_OFFICIAL,
			}
		case "old":
			val.Use = &r4dt.Identifier_UseCode{
				Value: r4gp.IdentifierUseCode_OLD,
			}
		case "secondary":
			val.Use = &r4dt.Identifier_UseCode{
				Value: r4gp.IdentifierUseCode_SECONDARY,
			}
		default:
			return fmt.Errorf("unknown use '%s'", use)
		}
		resource.Identifier = append(resource.Identifier, val)
		return nil
	}
}

func WithName(text, family string, given []string) WithFunc {
	return func(resource *r4pprac.Practitioner) error {
		if resource.Name == nil {
			resource.Name = make([]*r4dt.HumanName, 0)
		}
		var givenList []*r4dt.String
		for _, g := range given {
			givenList = append(givenList, &r4dt.String{Value: g})
		}
		resource.Name = append(resource.Name, &r4dt.HumanName{
			Text:   &r4dt.String{Value: text},
			Given:  givenList,
			Family: &r4dt.String{Value: family},
		})
		return nil
	}
}

func NewPractitioner(options ...WithFunc) (*r4pprac.Practitioner, error) {
	resource := &r4pprac.Practitioner{}

	for _, w := range options {
		if err := w(resource); err != nil {
			return nil, err
		}
	}
	return resource, nil
}
