package iron_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/stretchr/testify/assert"
)

var (
	muxIRON    *http.ServeMux
	serverIRON *httptest.Server
	client     *iron.Client
	projectID  = "48a0183d-a588-41c2-9979-737d15e9e860"
	token      = "YM7eZakYwqoui5znoH4g"
)

func setup(t *testing.T) func() {
	muxIRON = http.NewServeMux()
	serverIRON = httptest.NewServer(muxIRON)

	var err error

	client, err = iron.NewClient(&iron.Config{
		BaseURL:   serverIRON.URL,
		ProjectID: projectID,
		Token:     token,
		Debug:     true,
		DebugLog:  "/tmp/iron_test.log",
	})
	assert.Nil(t, err)
	assert.NotNil(t, client)

	return func() {
		serverIRON.Close()
	}
}

func TestClient_Error(t *testing.T) {
	muxIRON = http.NewServeMux()
	serverIRON = httptest.NewServer(muxIRON)

	var err error
	client, err = iron.NewClient(&iron.Config{
		BaseURL:   serverIRON.URL,
		ProjectID: projectID,
		Token:     token,
		Debug:     true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, client)

	task, resp, err := client.Tasks.GetTasks()
	assert.NotNil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NotNil(t, task)
	client.Close()
}
