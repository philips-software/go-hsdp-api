package s3creds

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func testPolicy() *Policy {
	p := &Policy{}
	p.Conditions.ManagingOrganizations = []string{"foo", "bar"}
	p.Conditions.Groups = []string{"low", "high"}
	p.Allowed.Resources = []string{"one", "two"}
	p.Allowed.Actions = []string{"GET", "PUT"}
	return p
}

func TestPolicyEqual(t *testing.T) {
	a := testPolicy()
	p := testPolicy()

	assert.True(t, a.Equals(p))

	p.Allowed.Actions = []string{"PUT", "GET"}
	assert.True(t, a.Equals(p))
	p.Allowed.Actions = []string{"POST", "DELETE"}
	assert.False(t, a.Equals(p))

	p = testPolicy()
	p.ID = 1
	assert.False(t, a.Equals(p))

	p = testPolicy()
	p.ResourceType = "bla"
	assert.False(t, a.Equals(p))

	p = testPolicy()
	p.Conditions.ManagingOrganizations = []string{"foo", "bar", "baz"}
	assert.False(t, a.Equals(p))

	p = testPolicy()
	a.Conditions.ManagingOrganizations = []string{"foo", "bar", "baz"}
	assert.False(t, a.Equals(p))
	a = testPolicy()

	p.Conditions.Groups = []string{"low"}
	assert.False(t, a.Equals(p))

	p = testPolicy()
	a.Conditions.Groups = []string{"low"}
	assert.False(t, a.Equals(p))
	a = testPolicy()

	p.Allowed.Resources = []string{"one", "two", "three"}
	assert.False(t, a.Equals(p))

	p = testPolicy()
	a.Allowed.Resources = []string{"one", "two", "three"}
	assert.False(t, a.Equals(p))
	a = testPolicy()

	p.Allowed.Actions = []string{"GET", "PUT", "DELETE"}
	assert.False(t, a.Equals(p))

	p = testPolicy()
	a.Allowed.Actions = []string{"GET", "PUT", "DELETE"}
	assert.False(t, a.Equals(p))
	a = testPolicy()

	p.ProductKey = "nokey"
	assert.False(t, a.Equals(p))
}

func TestValidatePolicy(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("policyActions", validateActions)
	assert.Nil(t, err)

	p := &Policy{}
	p.Conditions.ManagingOrganizations = []string{"foo", "bar"}
	p.Allowed.Resources = []string{"one", "two"}
	p.Allowed.Actions = []string{"GET", "PUT", "LIST", "DELETE", "ALL_OBJECT"}
	err = v.Struct(p)
	assert.Nilf(t, err, "Validation failed: %v", err)

	p.Allowed.Actions = []string{"POST"}
	err = v.Struct(p)
	assert.NotNil(t, err)
}
