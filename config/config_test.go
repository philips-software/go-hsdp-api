package config_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/philips-software/go-hsdp-api/config"
	"github.com/stretchr/testify/assert"
)

func localConfig(t *testing.T) (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !assert.True(t, ok) {
		return nil, fmt.Errorf("runtime.Caller(0) error")
	}
	basePath := filepath.Dir(filename)
	hsdpTomlFile := filepath.Join(basePath, "hsdp.toml")
	data, err := ioutil.ReadFile(hsdpTomlFile)
	if !assert.Nil(t, err) {
		return nil, err
	}
	configReader := bytes.NewReader(data)
	return config.New(config.FromReader(configReader))
}

func TestNew(t *testing.T) {
	c, err := localConfig(t)
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
	c, err := localConfig(t)
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
	_, filename, _, ok := runtime.Caller(0)
	if !assert.True(t, ok) {
		return
	}
	basePath := filepath.Dir(filename)
	hsdpTomlFile := filepath.Join(basePath, "hsdp.toml")
	data, err := ioutil.ReadFile(hsdpTomlFile)
	if !assert.Nil(t, err) {
		return
	}
	configReader := bytes.NewReader(data)
	c, err := config.New(
		config.WithEnv("client-test"),
		config.WithRegion("us-east"),
		config.FromReader(configReader))
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
	c, err := localConfig(t)
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

func TestRegions(t *testing.T) {
	c, err := localConfig(t)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	regions := c.Regions()
	assert.Less(t, 0, len(regions))
	assert.Contains(t, regions, "eu-west")
}

func TestServices(t *testing.T) {
	c, err := localConfig(t)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	services := c.Region("us-east").Env("client-test").Services()
	assert.Less(t, 0, len(services))
}

func TestKeys(t *testing.T) {
	c, err := localConfig(t)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	cartel := c.Region("us-east").Env("client-test").Service("cartel")
	assert.True(t, cartel.Available())
	keys := cartel.Keys()
	assert.Less(t, 0, len(keys))
	_, err = cartel.GetString("bogus")
	assert.NotNil(t, err)
	port, err := cartel.GetInt("port")
	assert.Equal(t, config.ErrNotFound, err)
	assert.Equal(t, 0, port)
}
