package logging

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func shouldRun(t *testing.T) bool {
	key := os.Getenv("INT_LOGGING_SHARED_KEY")
	secret := os.Getenv("INT_LOGGING_SECRET_KEY")
	productKey := os.Getenv("INT_LOGGING_PRODUCT_KEY")
	ingestorURL := os.Getenv("INT_LOGGING_INGESTOR_URL")

	if key == "" || secret == "" || productKey == "" || ingestorURL == "" {
		t.Skip("skipping integration test")
		return false
	}
	return true
}

func TestIntegration(t *testing.T) {
	if !shouldRun(t) {
		return
	}

	key := os.Getenv("INT_LOGGING_SHARED_KEY")
	secret := os.Getenv("INT_LOGGING_SECRET_KEY")
	productKey := os.Getenv("INT_LOGGING_PRODUCT_KEY")
	ingestorURL := os.Getenv("INT_LOGGING_INGESTOR_URL")

	if !assert.NotEmpty(t, key) {
		return
	}
	if !assert.NotEmpty(t, secret) {
		return
	}
	if !assert.NotEmpty(t, productKey) {
		return
	}
	if !assert.NotEmpty(t, ingestorURL) {
		return
	}

	intClient, err := NewClient(nil, &Config{
		SharedKey:    key,
		SharedSecret: secret,
		ProductKey:   productKey,
		BaseURL:      ingestorURL,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, intClient) {
		return
	}

	// Happy flow
	resp, err := intClient.StoreResources([]Resource{validResource}, 1)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	// Local validation test
	resp, err = intClient.StoreResources([]Resource{
		validResource,
		invalidResource,
		validResource,
		validResource,
		invalidResource,
		validResource,
		invalidResource}, 7)
	if !assert.NotNil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, ErrBatchErrors, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	assert.Equal(t, 3, len(resp.Failed))
}
