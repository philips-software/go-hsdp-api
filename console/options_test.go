package console_test

import (
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/console"
	"github.com/stretchr/testify/assert"
)

func TestWithStart(t *testing.T) {
	var req http.Request

	opt := console.WithStart(100)
	err := opt(&req)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, req.URL.RawQuery, "start=100")
}

func TestWithEnd(t *testing.T) {
	var req http.Request

	opt := console.WithEnd(424242)
	err := opt(&req)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, req.URL.RawQuery, "end=424242")
}

func TestWithStep(t *testing.T) {
	var req http.Request

	opt := console.WithStep(1234)
	err := opt(&req)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, req.URL.RawQuery, "step=1234")
}

func TestWithHost(t *testing.T) {
	var req http.Request

	opt := console.WithHost("abc.k8s.io")
	err := opt(&req)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, req.URL.RawQuery, "host=abc.k8s.io")
}

func TestWithQuery(t *testing.T) {
	var req http.Request

	opt := console.WithQuery("some_prom_metric{}")
	err := opt(&req)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, req.URL.RawQuery, "query=some_prom_metric%7B%7D")
}
