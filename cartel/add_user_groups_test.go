package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUserGroups(t *testing.T) {
	var addResponse = `{"foo.dev.com": {"ldap_groups": ["foo", "bar"]}}`

	teardown, err := setup(t, &Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})

	muxCartel.HandleFunc("/v3/api/add_ldap_group", endpointMocker([]byte(sharedSecret),
		addResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	aur, resp, err := client.AddUserGroups([]string{"foo.dev"}, []string{"foo", "bar"})

	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, aur) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, true, aur.Success())
}
