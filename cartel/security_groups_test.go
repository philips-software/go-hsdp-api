package cartel

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestSecurityGroups(t *testing.T) {
	teardown, err := setup(t, Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/get_security_groups", endpointMocker(sharedSecret,
		`[
    "foo",
    "bar",
    "baz"
]`))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	groups, resp, err := client.GetSecurityGroups()
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, groups) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, len(*groups))
}
