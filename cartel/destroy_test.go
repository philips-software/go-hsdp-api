package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestroy(t *testing.T) {
	var stopResponse = `{"message": {"foo.dev": {"cartel": "Instance suspended"}}}`
	var _ = `{"message": "Instance cannot be started due to current state: running"}`

	teardown, err := setup(t, Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/destroy", endpointMocker(sharedSecret,
		stopResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	sr, resp, err := client.Destroy("foo.dev")
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sr) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, sr.Success())
}
