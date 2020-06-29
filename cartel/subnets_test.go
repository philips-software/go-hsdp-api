package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubnets(t *testing.T) {
	var subnetResponse = `{
  "sc1-priv-e": {
    "id": "subnet-dead",
    "network": "192.68.24.0/21"
  },
  "sc1-private-a": {
    "id": "subnet-beaf",
    "network": "192.68.0.0/21"
  },
  "sc1-private-c": {
    "id": "subnet-treaty0self",
    "network": "192.68.8.0/21"
  }
}`
	teardown, err := setup(t, Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/get_all_subnets", endpointMocker([]byte(sharedSecret),
		subnetResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	subnets, resp, err := client.GetAllSubnets()
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, subnets) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, len(*subnets))
}
