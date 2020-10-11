package cartel

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBody(t *testing.T) {
	var requestBody RequestBody

	err := json.Unmarshal([]byte(`{}`), &requestBody)
	assert.Nil(t, err)
	m := requestBody.ToJson()
	assert.Equal(t, "{\"encrypt_vols\":false,\"protect\":false}", string(m))

	err = InstanceType("container-host")(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "container-host", requestBody.InstanceType)

	err = UserGroups("foo", "bar")(&requestBody)
	assert.Nil(t, err)
	assert.Contains(t, requestBody.LDAPGroups, "foo")
	assert.Contains(t, requestBody.LDAPGroups, "bar")

	err = VolumeType("gp2")(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "gp2", requestBody.VolumeType)

	err = IOPs(5000)(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, 5000, requestBody.IOPs)

	err = InSubnet("subnet")(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "subnet", requestBody.Subnet)

	err = SubnetType("private")(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "private", requestBody.SubnetType)

	err = VPCID("v-bbiab")(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "v-bbiab", requestBody.VpcId)

	err = Protect(true)(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, true, requestBody.Protect)

	err = VolumeEncryption(true)(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, true, requestBody.EncryptVols)

	err = Tags(map[string]string{
		"foo": "bar",
		"bar": "baz",
	})(&requestBody)
	assert.Nil(t, err)
	assert.Equal(t, "bar", requestBody.Tags["foo"])
	assert.Equal(t, "baz", requestBody.Tags["bar"])
}
