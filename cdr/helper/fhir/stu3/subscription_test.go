package stu3_test

import (
	"testing"

	"github.com/google/fhir/go/jsonformat"

	"github.com/stretchr/testify/assert"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
)

func TestNewSubscription(t *testing.T) {
	sub, err := stu3.NewSubscription(
		stu3.WithContact("phone", "(603) 203-2594", "work"),
		stu3.WithCriteria("Patient?given=Ron"),
		stu3.WithEndpoint("https://foo/notification"),
		stu3.WithHeaders([]string{"Authorization: Bearer xxx"}),
		stu3.WithReason("somereason"))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, sub) {
		return
	}
	ma, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, ma) {
		return
	}
	_, err = ma.MarshalResource(sub)
	if !assert.Nil(t, err) {
		return
	}
}
