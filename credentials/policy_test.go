package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyEqual(t *testing.T) {
	a := &Policy{}
	a.Conditions.ManagingOrganizations = []string{"foo", "bar"}
	a.Conditions.Groups = []string{"low", "high"}
	a.Allowed.Resources = []string{"one", "two"}
	a.Allowed.Actions = []string{"GET", "PUT"}
	b := &Policy{}
	b.Conditions.ManagingOrganizations = []string{"foo", "bar"}
	b.Conditions.Groups = []string{"low", "high"}
	b.Allowed.Resources = []string{"one", "two"}
	b.Allowed.Actions = []string{"GET", "PUT"}

	equal := a.Equals(b)
	assert.Equal(t, true, equal)

	b.Allowed.Actions = []string{"POST", "DELETE"}
	equal = a.Equals(b)
	assert.Equal(t, false, equal)
}
