package cartel

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(t, "gw-eu1.phsdp.com", bastion)
}
