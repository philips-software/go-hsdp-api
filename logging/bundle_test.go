package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundle(t *testing.T) {
	var b = Bundle{
		ResourceType: "transaction",
	}
	assert.Equal(t, b.ResourceType, "transaction")
}
