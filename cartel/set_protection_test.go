package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetProtect(t *testing.T) {
	var protectResponse = `{"message": "Termination protection enabled for: foo.dev.com"}`

	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/protect", endpointMocker([]byte(sharedSecret),
		protectResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	pr, resp, err := client.SetProtection("foo", true)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, pr) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, pr.Success())
}
