package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	teardown, err := setup(t, Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})

	muxCartel.HandleFunc("/v3/api/get_all_roles", endpointMocker(sharedSecret,
		`[
    {
        "description": "Some role.",
        "role": "amp"
    },
    {
        "description": "Another role.",
        "role": "base"
    },
    {
        "description": "And another one.",
        "role": "djkhaled"
    }]`))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	roles, resp, err := client.GetRoles()
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, roles) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Equal(t, 3, len(*roles)) {
		return
	}
	assert.Equal(t, "djkhaled", (*roles)[2].Role)
}
