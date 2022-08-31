package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTags(t *testing.T) {
	var addResponse = `{"message": "Added tags foo, bar to foo.dev"}`

	teardown, err := setup(t, &Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})

	muxCartel.HandleFunc("/v3/api/add_tags", endpointMocker([]byte(sharedSecret),
		addResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	atr, resp, err := client.AddTags([]string{"foo.dev"}, map[string]string{
		"foo": "bar",
		"bar": "baz",
	})
	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, atr) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, true, atr.Success())
}
