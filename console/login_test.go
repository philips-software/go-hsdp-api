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
	tk, err := client.SetToken(token).Token()
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, token, tk.AccessToken)
	client.SetTokens(token, token2, token, time.Now().Add(10*time.Minute).Unix())
	assert.Equal(t, token2, client.RefreshToken())
}
