package console_test

import (
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	token2 := "TokenMy"
	assert.Equal(t, token, client.SetToken(token).Token())
	client.SetTokens(token, token2, token, time.Now().Add(10*time.Minute).Unix())
	assert.Equal(t, token2, client.RefreshToken())
}
