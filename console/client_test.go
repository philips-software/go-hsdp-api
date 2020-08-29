package console

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	errors "golang.org/x/xerrors"
)

var (
	muxUAA        *http.ServeMux
	serverUAA     *httptest.Server
	muxCONSOLE    *http.ServeMux
	serverCONSOLE *httptest.Server
	token         string
	refreshToken  string

	client *Client
)

func setup(t *testing.T) (func(), error) {
	muxUAA = http.NewServeMux()
	serverUAA = httptest.NewServer(muxUAA)
	muxCONSOLE = http.NewServeMux()
	serverCONSOLE = httptest.NewServer(muxCONSOLE)
	var err error

	assert.Nil(t, err)

	client, err = NewClient(nil, &Config{
		UAAURL:         serverUAA.URL,
		BaseConsoleURL: serverCONSOLE.URL,
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid client")
		return func() {}, err
	}

	token = "44d20214-7879-4e35-923d-f9d4e01c9746"
	token2 := "55d20214-7879-4e35-923d-f9d4e01c9746"
	refreshToken = "31f1a449-ef8e-4bfc-a227-4f2353fde547"

	muxUAA.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
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

	return func() {
		serverUAA.Close()
		serverCONSOLE.Close()
	}, nil
}

func TestLogin(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	err = client.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token, client.Token())
	assert.Equal(t, refreshToken, client.RefreshToken())
}

func TestConsoleRequest(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	req, err := client.NewRequest(CONSOLE, "GET", "/foo", nil, nil)
	if err != nil {
		t.Errorf("Expected no no errors, got: %v", err)
	}
	if req == nil {
		t.Errorf("Expected valid request")
	}
	req, _ = client.NewRequest(CONSOLE, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			r.Header.Set("Foo", "Bar")
			return nil
		},
	})
	if req.Header.Get("Foo") != "Bar" {
		t.Errorf("Expected OptionFuncs to be processed")
	}
	testErr := errors.New("test error")
	req, err = client.NewRequest(CONSOLE, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			return testErr
		},
	})
	assert.Nil(t, req)
	assert.Equal(t, testErr, err)
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

	client, err = NewClient(nil, &Config{
		UAAURL:         serverUAA.URL,
		BaseConsoleURL: serverCONSOLE.URL,
		Debug:          true,
		DebugLog:       tmpfile.Name(),
	})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer client.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	err = client.Login("username", "password")
	assert.Nil(t, err)

	fi, err := tmpfile.Stat()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("Expected something to be written to DebugLog")
	}

}

func TestTokenRefresh(t *testing.T) {
	muxUAA = http.NewServeMux()
	serverUAA = httptest.NewServer(muxUAA)
	muxCONSOLE = http.NewServeMux()
	serverCONSOLE = httptest.NewServer(muxCONSOLE)

	defer serverUAA.Close()
	defer serverCONSOLE.Close()

	cfg := &Config{
		UAAURL:         serverUAA.URL,
		BaseConsoleURL: serverCONSOLE.URL,
		Scopes:         []string{"introspect", "cn"},
	}
	client, err := NewClient(nil, cfg)
	assert.Nil(t, err)

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"
	refreshToken := "13614f90-9cdf-4962-aea3-01cd51fa56b9"
	newToken := "90b208cd-aaf3-45bb-9410-ba3f42255b9d"
	newRefreshToken := "9c45339e-38c8-4dac-b290-5c3ac571c369"

	muxUAA.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "POST", r.Method) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := r.ParseForm(); !assert.Nilf(t, err, "Unable to parse form") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		grantType := strings.Join(r.Form["grant_type"], " ")
		receveidRefreshToken := strings.Join(r.Form["refresh_token"], " ")

		w.Header().Set("Content-Type", "application/json")
		switch grantType {
		case "refresh_token":
			if !assert.Equal(t, refreshToken, receveidRefreshToken) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"scope": "`+strings.Join(cfg.Scopes, " ")+`",
				"access_token": "`+newToken+`",
				"refresh_token": "`+newRefreshToken+`",
				"expires_in": 1799,
				"token_type": "Bearer"
			}`)
		case "password":
			err := r.ParseForm()
			assert.Nil(t, err)
			username := r.Form.Get("username")
			if !assert.Equal(t, "username", username) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"scope": "`+strings.Join(cfg.Scopes, " ")+`",
				"access_token": "`+token+`",
				"refresh_token": "`+refreshToken+`",
				"expires_in": 1799,
				"token_type": "Bearer"
			}`)
		}

	})
	err = client.Login("username", "password")
	assert.Nil(t, err)

	err = client.tokenRefresh()
	assert.Nilf(t, err, fmt.Sprintf("Unexpected error: %v", err))
	assert.Equal(t, newToken, client.Token())
	assert.Equal(t, newRefreshToken, client.RefreshToken())
	httpClient := client.HttpClient()
	assert.NotNil(t, httpClient)
}

func TestAutoconfig(t *testing.T) {
	cfg := &Config{
		Region:      "us-east",
		Environment: "client-test",
	}
	// Explicit config always wins over autoconfig
	foo := "https://foo.com"
	cfg.BaseConsoleURL = foo
	cfg.UAAURL = foo
	_, _ = NewClient(nil, cfg)
	assert.Equal(t, foo, cfg.BaseConsoleURL)
	assert.Equal(t, foo, cfg.UAAURL)
}
