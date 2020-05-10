package cartel

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	muxCartel    *http.ServeMux
	cartelServer *httptest.Server
	client       *Client
)

var (
	sharedSecret = []byte("SharedSecret")
	sharedToken  = "SharedToken"
)

func endpointMocker(secret []byte, responseBody string) func(http.ResponseWriter, *http.Request) {
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}
}

func setup(t *testing.T, config Config) (func(), error) {
	var err error

	muxCartel = http.NewServeMux()
	cartelServer = httptest.NewServer(muxCartel)
	if config.Host != "" { // So we can test for missing BaseURL
		config.Host = cartelServer.URL
	}

	client, err = NewClient(nil, config)
	if err != nil {
		return func() {
			cartelServer.Close()
		}, err
	}

	muxCartel.HandleFunc("/v3/api/get_security_groups", endpointMocker(config.Secret,
		`[
    "foo",
    "bar",
    "baz"
]`))

	return func() {
		cartelServer.Close()
	}, nil
}

func TestGetRoles(t *testing.T) {
	teardown, err := setup(t, Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "http://foo",
	})
	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.GetSecurityGroups()
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
