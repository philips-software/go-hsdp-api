package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestroy(t *testing.T) {
	var destroyResponse = `{
  "AWS": "Instance(s) i-xxa562c262b18bfzz terminated",
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

	muxCartel.HandleFunc("/v3/api/destroy", endpointMocker([]byte(sharedSecret),
		destroyResponse))

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
