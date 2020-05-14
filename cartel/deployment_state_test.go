package cartel

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentState(t *testing.T) {
	var deployResponse = `{"foo.dev":{"deploy_state":"succeeded"}}`

	teardown, err := setup(t, Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
	})

	muxCartel.HandleFunc("/v3/api/deployment_status", endpointMocker(sharedSecret,
		deployResponse))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	aur, resp, err := client.GetDeploymentState("foo.dev")

	assert.Nil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "succeeded", aur)
}
