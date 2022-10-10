package stu3_test

import (
	"testing"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"

	"github.com/stretchr/testify/assert"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
)

func TestNewSubscription(t *testing.T) {
	endTime := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)

	deleteEndpoint := "https://foo/delete_notification"

	sub, err := stu3.NewSubscription(
		stu3.WithContact("phone", "(603) 203-2594", "work"),
		stu3.WithCriteria("Patient?given=Ron"),
		stu3.WithEndpoint("https://foo/notification"),
		stu3.WithDeleteEndpoint(deleteEndpoint),
		stu3.WithHeaders([]string{"Authorization: Bearer cm9uOnN3YW5zb24="}),
		stu3.WithReason("some reason"),
		stu3.WithEndtime(endTime))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, sub) {
		return
	}
	assert.Equal(t, endTime.UnixNano()/1000, sub.End.ValueUs)
	ma, err := jsonformat.NewMarshaller(false, "", "", fhirversion.STU3)
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
	getDEV := stu3.DeleteEndpointValue()
	assert.Equal(t, deleteEndpoint, getDEV(sub))
}
