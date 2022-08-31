package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllInstances(t *testing.T) {
	teardown, err := setup(t, &Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})
	var responseBody = `[{"instance_id":"i-deadbeaf","name_tag":"some.dev","owner":"xxx","role":"container-host"}]`

	muxCartel.HandleFunc("/v3/api/get_all_instances", endpointMocker([]byte(sharedSecret),
		responseBody))
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	instances, resp, err := client.GetAllInstances()
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, instances) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	if !assert.Equal(t, 1, len(*instances)) {
		return
	}
	assert.Equal(t, "xxx", (*instances)[0].Owner)
}
