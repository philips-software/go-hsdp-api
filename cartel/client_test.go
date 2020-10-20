package cartel

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	muxCartel    *http.ServeMux
	cartelServer *httptest.Server
	client       *Client
)

var (
	sharedSecret = "SharedSecret"
	sharedToken  = "SharedToken"
)

func endpointMocker(secret []byte, responseBody string, statusCode ...int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		signature := generateSignature(secret, body)
		auth := r.Header.Get("Authorization")
		if signature != auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if len(statusCode) > 0 {
			w.WriteHeader(statusCode[0])
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_, _ = w.Write([]byte(responseBody))
	}
}

func setup(t *testing.T, config *Config) (func(), error) {
	var err error

	muxCartel = http.NewServeMux()
	cartelServer = httptest.NewServer(muxCartel)
	if config.Host != "" { // So we can test for missing BaseURL
		u, _ := url.Parse(cartelServer.URL)
		config.Host = u.Host
	}

	client, err = NewClient(nil, config)
	if err != nil {
		return func() {
			cartelServer.Close()
		}, err
	}

	return func() {
		cartelServer.Close()
	}, nil
}

func TestDebug(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "cartel")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	teardown, err := setup(t, &Config{
		NoTLS:      true,
		SkipVerify: true,
		Token:      sharedToken,
		Secret:     sharedSecret,
		Host:       "foo",
		Debug:      true,
		DebugLog:   tmpfile.Name(),
	})
	var responseBody = `[{"instance_id":"i-deadbeaf","name_tag":"some.dev","owner":"xxx","role":"container-host"}]`

	muxCartel.HandleFunc("/v3/api/get_all_instances", endpointMocker([]byte(sharedSecret),
		responseBody))
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	_, _, err = client.GetAllInstances()
	if !assert.Nil(t, err) {
		return
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up
	fi, err := tmpfile.Stat()
	if !assert.Nil(t, err) {
		return
	}
	assert.Less(t, int64(0), fi.Size())
}

func TestAutoconfig(t *testing.T) {
	cfg := &Config{
		Token:  "alice",
		Secret: "foo",
		Region: "us-east",
	}
	_, err := NewClient(nil, cfg)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, cfg.Host)

	// Explicit config always wins over autoconfig
	foo := "foo.com"
	cfg.Host = foo
	_, _ = NewClient(nil, cfg)
	assert.Equal(t, foo, cfg.Host)
}
