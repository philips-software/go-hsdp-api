package cartel

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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

func setup(t *testing.T, config Config) (func(), error) {
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
