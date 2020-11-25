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
	hsdpJsonFile := filepath.Join(basePath, "hsdp.json")
	data, err := ioutil.ReadFile(hsdpJsonFile)
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
	assert.Equal(t, "https://iam-client-test.us-east.philips-healthsuite.com", iamService.URL)
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
	assert.Equal(t, "cartel-na1.cloud.phsdp.com", cartelService.Host)
}

func TestOpts(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !assert.True(t, ok) {
		return
	}
	basePath := filepath.Dir(filename)
	hsdpJsonFile := filepath.Join(basePath, "hsdp.json")
	data, err := ioutil.ReadFile(hsdpJsonFile)
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
	cartel := c.Service("cartel")
	assert.Equal(t, "cartel-na1.cloud.phsdp.com", cartel.Host)
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
	assert.Equal(t, "", missingService.URL)
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
