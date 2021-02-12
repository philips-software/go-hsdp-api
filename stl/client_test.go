package stl_test

import (
	"context"
	"github.com/hasura/go-graphql-client"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	muxUAA        *http.ServeMux
	serverUAA     *httptest.Server
	muxCONSOLE    *http.ServeMux
	serverCONSOLE *httptest.Server
	muxSTL        *http.ServeMux
	serverSTL     *httptest.Server
	token         string
	refreshToken  string

	consoleClient *console.Client
	client        *stl.Client
)

func setup(t *testing.T) (func(), error) {
	muxUAA = http.NewServeMux()
	serverUAA = httptest.NewServer(muxUAA)
	muxCONSOLE = http.NewServeMux()
	serverCONSOLE = httptest.NewServer(muxCONSOLE)
	muxSTL = http.NewServeMux()
	serverSTL = httptest.NewServer(muxSTL)
	var err error

	assert.Nil(t, err)

	consoleClient, err = console.NewClient(nil, &console.Config{
		UAAURL:         serverUAA.URL,
		BaseConsoleURL: serverCONSOLE.URL,
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
	client, err = stl.NewClient(consoleClient, &stl.Config{
		STLAPIURL: serverSTL.URL,
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid STL client")
		return func() {}, err
	}

	return func() {
		serverUAA.Close()
		serverCONSOLE.Close()
		serverSTL.Close()
	}, nil
}

func TestDebug(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	client, err = stl.NewClient(consoleClient, &stl.Config{
		STLAPIURL: serverSTL.URL,
		DebugLog:  tmpfile.Name(),
	})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer client.Close()
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up

	var query struct {
		App stl.AppResource `graphql:"applicationResource(id: $id, name: $name)"`
	}
	err = client.Query(context.Background(), &query, map[string]interface{}{
		"id":   graphql.Int(1),
		"name": graphql.String("name"),
	})
	assert.NotNil(t, err)

	fi, err := tmpfile.Stat()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("Expected something to be written to DebugLog")
	}
}
