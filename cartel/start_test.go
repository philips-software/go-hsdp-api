package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	var startResponse = `{"message": {"foo.dev": {"cartel": "Instance started"}}}`
	var _ = `{"message": "Instance cannot be started due to current state: running"}`

	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/start", endpointMocker([]byte(sharedSecret),
		startResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	sr, resp, err := client.Start("foo.dev")
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

func TestAlreadyRunning(t *testing.T) {
	var failResponse = `{"message": "Instance cannot be started due to current state: running"}`

	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/start", endpointMocker([]byte(sharedSecret),
		failResponse, http.StatusBadRequest))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	sr, resp, err := client.Start("foo.dev")
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sr) {
		return
	}
	if !assert.NotNil(t, err) {
		return
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, false, sr.Success())
}
