package iam_test

import (
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfiles(t *testing.T) {
	a := iam.Profile{
		Addresses: []iam.Address{
			{},
			{
				IsPrimary: "yes",
			},
			{
				Street:    "",
				Building:  "",
				IsPrimary: "",
			},
			{},
		},
	}
	assert.Equal(t, 4, len(a.Addresses))
	a.PruneBlankAddresses()
	assert.Equal(t, 1, len(a.Addresses))
}
