package config_test

import (
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	iamService := c.
		Region("us-east").
		Env("client-test").
		Service("iam")
	if !assert.NotNil(t, iamService) {
		return
	}
	url, err := iamService.GetString("iam_url")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "https://iam-client-test.us-east.philips-healthsuite.com", url)
}

func TestCartel(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	cartelService := c.
		Region("us-east").
		Service("cartel")
	if !assert.NotNil(t, cartelService) {
		return
	}
	host, err := cartelService.GetString("host")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "cartel-na1.cloud.phsdp.com", host)
}

func TestOpts(t *testing.T) {
	resp, err := http.Get(config.CanonicalURL)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	defer resp.Body.Close()
	c, err := config.New(
		config.WithEnv("client-test"),
		config.WithRegion("us-east"),
		config.FromReader(resp.Body))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	host, err := c.Service("cartel").GetString("host")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "cartel-na1.cloud.phsdp.com", host)

}

func TestMissing(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	missingService := c.
		Region("us-east").
		Service("bogus")
	assert.False(t, missingService.Available())
	_, err = missingService.GetString("foo")
	assert.NotNil(t, err)
}

func TestServices(t *testing.T) {
	c, err := config.New(
		config.WithRegion("us-east"),
		config.WithEnv("client-test"))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	services := c.Services()
	assert.Less(t, 0, len(services))
}

func TestKeys(t *testing.T) {
	c, err := config.New(
		config.WithRegion("us-east"),
		config.WithEnv("client-test"))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	cartel := c.Service("cartel")
	assert.True(t, cartel.Available())
	keys := cartel.Keys()
	assert.Less(t, 0, len(keys))
	_, err = cartel.GetString("bogus")
	assert.NotNil(t, err)
	port, err := cartel.GetInt("port")
	assert.Equal(t, config.ErrNotFound, err)
	assert.Equal(t, 0, port)
}
