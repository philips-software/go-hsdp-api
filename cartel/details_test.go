package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetails(t *testing.T) {
	teardown, err := setup(t, Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})
	var responseBody = `[
  {
    "xxx.dev": {
      "block_devices": [
        "/dev/sdz",
        "/dev/sda1"
      ],
      "instance_id": "i-deadbeaf",
      "instance_type": "m5.large",
      "launch_time": "2019-11-29T18:17:44.000Z",
      "ldap_groups": [
        "my-group"
      ],
      "private_address": "192.168.99.66",
      "protection": false,
      "public_address": null,
      "role": "container-host",
      "security_groups": [
        "foo",
        "bar"
      ],
      "state": "running",
      "subnet": "some-private-network",
      "tags": {
        "billing": "{\"token\": \"stoomlocomotief\", \"username\": \"xxx\"}"
      },
      "vpc": "vpca-bc",
      "zone": "us-east-3a"
    }
  }
]`

	muxCartel.HandleFunc("/v3/api/instance_details", endpointMocker(sharedSecret,
		responseBody))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	details, resp, err := client.GetDetails("xxx.dev")
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, details) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "container-host", details.Role)
}
