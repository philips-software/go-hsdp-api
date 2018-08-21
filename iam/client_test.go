package iam

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	signer "github.com/philips-software/go-hsdp-signer"
)

var (
	muxIAM     *http.ServeMux
	serverIAM  *httptest.Server
	muxIDM     *http.ServeMux
	serverIDM  *httptest.Server
	signerHSDP *signer.Signer

	client *Client
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)
	sharedKey := "SharedKey"
	secretKey := "SecretKey"

	signerHSDP, _ = signer.New(sharedKey, secretKey)

	client, _ = NewClient(nil, &Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      sharedKey,
		SecretKey:      secretKey,
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
	})

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
    "scope": "auth_iam_introspect mail tdr.contract tdr.dataitem",
    "access_token": "`+token+`",
    "refresh_token": "31f1a449-ef8e-4bfc-a227-4f2353fde547",
    "expires_in": "1799",
    "token_type": "Bearer"
}`)
	})

	return func() {
		serverIAM.Close()
		serverIDM.Close()
	}
}

func TestLogin(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	err := client.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	if client.Token() != token {
		t.Errorf("Unexpected token found: %s, expected: %s", client.Token(), token)
	}
}

func TestHasScopes(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	err := client.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	if !client.HasScopes("mail") {
		t.Errorf("Expected mail scope to be there")
	}
	if !client.HasScopes("tdr.contract", "tdr.dataitem") {
		t.Errorf("Expected tdr scopes to be there")
	}
	if client.HasScopes("missing") {
		t.Errorf("Unexpected scope confirmation")
	}
	if client.HasScopes("mail", "bogus") {
		t.Errorf("Unexpected scope confirmation")
	}
}

func TestIAMRequest(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	req, err := client.NewRequest(IAM, "GET", "/foo", nil, nil)
	if err != nil {
		t.Errorf("Expected no no errors, got: %v", err)
	}
	if req == nil {
		t.Errorf("Expected valid request")
	}
	req, _ = client.NewRequest(IAM, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			r.Header.Set("Foo", "Bar")
			return nil
		},
	})
	if req.Header.Get("Foo") != "Bar" {
		t.Errorf("Expected OptionFuncs to be processed")
	}
	testErr := errors.New("test error")
	req, err = client.NewRequest(IAM, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			return testErr
		},
	})
	if err == nil {
		t.Errorf("Request IAM request to fail")
	}
}

func TestIDMRequest(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	client.SetToken("xxx")
	req, err := client.NewRequest(IDM, "GET", "/foo", nil, nil)
	if err != nil {
		t.Errorf("Expected no no errors, got: %v", err)
	}
	if req == nil {
		t.Errorf("Expected valid request")
	}
	req, _ = client.NewRequest(IDM, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			r.Header.Set("Foo", "Bar")
			return nil
		},
	})
	if req.Header.Get("Foo") != "Bar" {
		t.Errorf("Expected OptionFuncs to be processed")
	}
	testErr := errors.New("test error")
	req, err = client.NewRequest(IDM, "POST", "/foo", nil, []OptionFunc{
		func(r *http.Request) error {
			return testErr
		},
	})
	if err == nil {
		t.Errorf("Request IDM request to fail")
	}
}
