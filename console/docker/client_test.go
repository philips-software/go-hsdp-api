package docker_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/console/docker"
	"github.com/stretchr/testify/assert"
)

var (
	muxUAA        *http.ServeMux
	serverUAA     *httptest.Server
	muxCONSOLE    *http.ServeMux
	serverCONSOLE *httptest.Server
	muxSTL        *http.ServeMux
	serverDocker  *httptest.Server
	token         string
	refreshToken  string

	consoleClient *console.Client
	client        *docker.Client
	tmpFile       *os.File
)

func setup(t *testing.T) (func(), error) {
	muxUAA = http.NewServeMux()
	serverUAA = httptest.NewServer(muxUAA)
	muxCONSOLE = http.NewServeMux()
	serverCONSOLE = httptest.NewServer(muxCONSOLE)
	muxSTL = http.NewServeMux()
	serverDocker = httptest.NewServer(muxSTL)
	var err error

	assert.Nil(t, err)

	tmpFile, err = ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	consoleClient, err = console.NewClient(nil, &console.Config{
		UAAURL:         serverUAA.URL,
		BaseConsoleURL: serverCONSOLE.URL,
		DebugLog:       tmpFile.Name(),
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid consoleClient")
		return func() {}, err
	}
	token = "44d20214-7879-4e35-923d-f9d4e01c9746"
	token2 := "55d20214-7879-4e35-923d-f9d4e01c9746"
	refreshToken = "31f1a449-ef8e-4bfc-a227-4f2353fde547"

	muxUAA.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			assert.Equal(t, "POST", r.Method)
		}
		err := r.ParseForm()
		assert.Nil(t, err)
		username := r.Form.Get("username")
		returnToken := token
		if username == "username2" {
			returnToken = token2
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
    		"scope": "auth_iam_introspect mail",
    		"access_token": "`+returnToken+`",
    		"refresh_token": "`+refreshToken+`",
    		"expires_in": 1799,
    		"token_type": "Bearer"
		}`)
	})
	err = consoleClient.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	client, err = docker.NewClient(consoleClient, &docker.Config{
		DockerAPIURL: serverDocker.URL,
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid Docker client")
		return func() {}, err
	}

	return func() {
		serverUAA.Close()
		serverCONSOLE.Close()
		serverDocker.Close()
		_ = os.Remove(tmpFile.Name())
	}, nil
}

func TestDebug(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	client, err = docker.NewClient(consoleClient, &docker.Config{
		DockerAPIURL: serverDocker.URL,
	})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer client.Close()

	var query struct {
		Resources []struct {
			docker.ServiceKeyNode
		} `graphql:"serviceKeys"`
	}
	err = client.Query(context.Background(), &query, nil)
	assert.NotNil(t, err)

	fi, err := tmpFile.Stat()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("Expected something to be written to DebugLog")
	}
}
