package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSecurityGroups(t *testing.T) {
	var addResponse = `{"message": "Security group(s) foo added to foo.dev.com"}`

	teardown, err := setup(t, Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})

	muxCartel.HandleFunc("/v3/api/add_security_groups", endpointMocker(sharedSecret,
		addResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	aur, resp, err := client.AddSecurityGroups([]string{"foo.dev"}, []string{"foo"})

	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, aur) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, aur.Success())
}
