package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStop(t *testing.T) {
	var stopResponse = `{
  "AWS": "Instance(s) i-03a562c262b18bf3d terminated",
  "Cartel": {
    "foo.dev.com": "Instance removed."
  }
}`

	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/suspend", endpointMocker([]byte(sharedSecret),
		stopResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	sr, resp, err := client.Stop("foo.dev")
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sr) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, true, sr.Success())
}
