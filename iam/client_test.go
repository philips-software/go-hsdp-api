package iam

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	muxIAM    *http.ServeMux
	serverIAM *httptest.Server
	muxIDM    *http.ServeMux
	serverIDM *httptest.Server

	client *Client
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)

	client, _ = NewClient(nil, &Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      "SharedKey",
		SecretKey:      "SecretKey",
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
}
