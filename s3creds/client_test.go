package s3creds

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

var (
	muxIAM      *http.ServeMux
	serverIAM   *httptest.Server
	muxIDM      *http.ServeMux
	serverIDM   *httptest.Server
	muxCreds    *http.ServeMux
	serverCreds *httptest.Server

	iamClient   *iam.Client
	credsClient *Client
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)
	muxCreds = http.NewServeMux()
	serverCreds = httptest.NewServer(muxCreds)

	var err error
	iamClient, err = iam.NewClient(nil, &iam.Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      "SharedKey",
		SecretKey:      "SecretKey",
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create iamCleitn: %v", err)
	}
	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
    "scope": "mail",
    "access_token": "`+token+`",
    "refresh_token": "31f1a449-ef8e-4bfc-a227-4f2353fde547",
    "expires_in": 1799,
    "token_type": "Bearer"
}`)
	})

	// Login immediately so we can create tdrClient
	err = iamClient.Login("username", "password")
	assert.Nil(t, err)

	credsClient, err = NewClient(iamClient, &Config{
		BaseURL: serverCreds.URL,
	})
	assert.Nilf(t, err, "failed to create credsClient: %v", err)

	return func() {
		serverIAM.Close()
		serverIDM.Close()
	}
}

func TestLogin(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	err := iamClient.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	accessToken, err := iamClient.Token()
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, token, accessToken)
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	credsClient, err = NewClient(iamClient, &Config{
		BaseURL:  serverCreds.URL,
		Debug:    true,
		DebugLog: tmpfile.Name(),
	})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer credsClient.Close()
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up

	err = iamClient.Login("username", "password")
	assert.Nil(t, err)

	id := 1
	productKey := "803505cd-79de-4441-88d7-6b110cd62b6d"
	_, _, _ = credsClient.Policy.GetPolicy(&GetPolicyOptions{
		ID:         &id,
		ProductKey: &productKey,
	})

	fi, err := tmpfile.Stat()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, fi.Size(), "Expected something to be written to DebugLog")
}
