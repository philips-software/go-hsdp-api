package r4_test

import (
	"testing"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"

	"github.com/stretchr/testify/assert"

	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4"
)

func TestNewR4Subscription(t *testing.T) {
	endTime := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)

	deleteEndpoint := "https://foo/delete_notification"

	sub, err := r4.NewSubscription(
		r4.WithContact("phone", "(603) 203-2594", "work"),
		r4.WithCriteria("Patient?given=Ron"),
		r4.WithEndpoint("https://foo/notification"),
		r4.WithDeleteEndpoint(deleteEndpoint),
		r4.WithHeaders([]string{"Authorization: Bearer cm9uOnN3YW5zb24="}),
		r4.WithReason("some reason"),
		r4.WithEndtime(endTime))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, sub) {
		return
	}
	assert.Equal(t, endTime.UnixNano()/1000, sub.End.ValueUs)
	ma, err := jsonformat.NewMarshaller(false, "", "", fhirversion.R4)
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
	getDEV := r4.DeleteEndpointValue()
	assert.Equal(t, deleteEndpoint, getDEV(sub))
}
