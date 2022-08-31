package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupDetails(t *testing.T) {
	var detailsResponse = `{
  "tcp-1080": [
    {
      "port_range": "1080-1080",
      "protocol": "tcp",
      "source": [
        "base",
        "100.0.100.0/16"
      ]
    }
  ]
}`
	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/security_group_details", endpointMocker([]byte(sharedSecret),
		detailsResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	details, resp, err := client.GetSecurityGroupDetails("tcp-1080")
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, details) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	if !assert.Equal(t, 1, len(*details)) {
		return
	}
	assert.Equal(t, "tcp", (*details)[0].Protocol)

}
