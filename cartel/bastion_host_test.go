package cartel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBastionHost(t *testing.T) {
	teardown, err := setup(t, &Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		NoTLS:  true,
		Region: "eu-west",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	bastion := client.BastionHost()
	assert.Equal(t, "gw.eu1.hsdp.io", bastion)
}
