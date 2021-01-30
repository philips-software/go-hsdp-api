package console_test

import (
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetToken(t *testing.T) {
	client, err := console.NewClient(nil, &console.Config{
		BaseConsoleURL: "http://localhost/console",
		UAAURL:         "http://localhost/uaa",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}
	token := "MyToken"
	assert.Equal(t, token, client.SetToken(token).Token())
}
