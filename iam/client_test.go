package iam

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	signer "github.com/philips-software/go-hsdp-signer"
)

var (
	muxIAM       *http.ServeMux
	serverIAM    *httptest.Server
	muxIDM       *http.ServeMux
	serverIDM    *httptest.Server
	signerHSDP   *signer.Signer
	token        string
	refreshToken string

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

	token = "44d20214-7879-4e35-923d-f9d4e01c9746"
	refreshToken = "31f1a449-ef8e-4bfc-a227-4f2353fde547"

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
    		"scope": "auth_iam_introspect mail tdr.contract tdr.dataitem",
    		"access_token": "`+token+`",
    		"refresh_token": "`+refreshToken+`",
    		"expires_in": 1799,
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

	if client.BaseIAMURL().String() != serverIAM.URL+"/" {
		t.Errorf("Unexpected BaseIAMURL: %s <-> %s", client.BaseIAMURL().String(), serverIAM.URL)
	}
	if client.BaseIDMURL().String() != serverIDM.URL+"/" {
		t.Errorf("Unexpected BaseIDMURL: %s <-> %s", client.BaseIDMURL().String(), serverIDM.URL)
	}
	if client.RefreshToken() != refreshToken {
		t.Errorf("Unexpected refresh token")
	}
}

func TestLoginWithScopes(t *testing.T) {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)

	defer serverIAM.Close()
	defer serverIDM.Close()

	sharedKey := "SharedKey"
	secretKey := "SecretKey"

	cfg := &Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      sharedKey,
		SecretKey:      secretKey,
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
		Scopes:         []string{"introspect", "cn"},
	}
	client, _ = NewClient(nil, cfg)

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Unable to parse form")
		}
		if strings.Join(r.Form["scope"], " ") != "introspect cn" {
			t.Fatalf("Expected scope to be `introspect cn` in test")
		}
		if strings.Join(r.Form["grant_type"], " ") != "password" {
			t.Fatalf("Exepcted grant_type to be `password` in test")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
    		"scope": "`+strings.Join(cfg.Scopes, " ")+`",
    		"access_token": "`+token+`",
    		"refresh_token": "31f1a449-ef8e-4bfc-a227-4f2353fde547",
    		"expires_in": 1799,
    		"token_type": "Bearer"
		}`)
	})

	err := client.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	if !client.HasScopes("introspect", "cn") {
		t.Errorf("Expected `introspect` and `cn` scope to be there")
	}
}

func TestServiceLogin(t *testing.T) {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)

	defer serverIAM.Close()
	defer serverIDM.Close()

	sharedKey := "SharedKey"
	secretKey := "SecretKey"

	cfg := &Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      sharedKey,
		SecretKey:      secretKey,
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
		Scopes:         []string{"introspect", "cn"},
	}
	service := &Service{
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----MIIEpAIBAAKCAQEAwM8LhQS4OB6e0xrMHE20NI/vWAwdgG3eoa50mlhlDwKQg0/sMYUKZBHkcit4rEQvgpXb36WtBhLAGC5gxLCBioRMfFG6c+DS9xyKXCexTTQZC1qBZlh1M7kq6oywnqfozBJ/9nAneOIkqA4NT9sy7jSMDuGFursL7p0iB1LrqEptBxm1zZKOw9GXUzqGTa+jdVj4DoviBtm6DCnQ61ucOEkl6DGvll5QBI693XIomqIbBICRHeMcTNoJ2GmKPYRITazKyk7FJc7Sn7E5T+ZB9StkX4BgiZjjZwnUKYpYGX65tipy3nrTzzZHfM83Rl91Gn1fbWIGxipxhp544XcjpwIDAQABAoIBAQC81YrYumiaPhMbenFRfyDxIc8uEp+KOxEClNQKnmxLqR1UHiCb10r3+zYcQ0sqnJVTdeYkQiUVf6O3iySnPp+AxFYMpBbSiuzTrJYt74oMrOuiXP/C9vvCrqXDlgsdOCIeTDgbanieQg3YsfqDrZFSDxDlOic5XRwwlKDRP3siFBuLrZZ/PtzNO4uMeBSKxGveghdQJiQ07XZoZ7bBsRU5lIV7tI+bpryA+xLu6C/LSRgtnwvHxnrppdvVOB4ZMc1IcigmCvMKCdUG1AQuxFH/v8ACuMuDubYRCJECtKlvQHNryA/uF8FYINjWDFiyoUmg2uu+Xtk13dY2uJoEycshAoGBAP1C1xs/9gaji6kwMP449IPH0pS24Nm0Rv3nxN6BPp0N+WpqCs9OAIFeNUZzScI6rMUtOBxEbddLGCkXCWnpvaPUoApx0192aW7mZ3mfhYQsaHohQHjkrzOCuRULaEbyET//w+7FvmxRHrjrEkuu8Tg06Uxt9394f5qVCrx6Eo2ZAoGBAMLk10htsuHx1KwYw7r6SLGaHhV60viUfKhiVXWUZiYCfgBtcr8Q/XBNFy/7e6Sm5o5xDBknFFcraMFKXpr9UMyOLLZLRqh8rxdRtoI3bGR1pNnx1BjkF6JmbDn8zqcibTet8nRXf7HoozolpbqL0QF8IXPzZpLd7o/+4RF4mYM/AoGBAOxtgpxoyIeIE/A9Ee+yQenoGFlGpH/4QTH1NR827rn1erryBedjfStIRFnhdKEC35kvTqts4lHTQ9nQLLSYRbZ033cArf/3bhPeugibeCxcvKgO9L4nVruytI/F13IrtxjU7xevuMYrsI+Wu7y1s3DyTD1Sh3OTjSRFMQGkwD85AoGANrfJOayS7JzY+Ph6+6QJhNOgXqd9VA1ccmopVDm19DX+6l/QN5StkzoRqIcSz8eMM7HJk8ZFD7RAVQRsS1eTt9qy8vtvex6GiiWG+EhXRl1BS295/QMNH6th92XjH0mrIFbWG5P1Zh3KtiibvyRCKgiP294ajmGA+Sy2RBF4CEECgYAKlx5gd2/1iCjOtnWRMjzSfhPUKe6W7FE3QwgoT7ZTA2xM4gLsSfC148G2EZMjDdKrHUdLC+hAQ4tRAqLMLjfyIhU/SKMM4NnVZlcZqdbWJP2BviXUz1rGKlwyt8fhvzbeiRbqTQxDdl0cjpE4Zi8aFSyczOGBsxQ0X57n5ARAtA==-----END RSA PRIVATE KEY-----",
		ServiceID:  "testservice.vhtestapp.vhprop@vhdev.h2hdevorg.philips-healthsuite.com",
	}

	client, _ = NewClient(nil, cfg)

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}
		r.ParseForm()
		if strings.Join(r.Form["grant_type"], " ") != "urn:ietf:params:oauth:grant-type:jwt-bearer" {
			t.Fatalf("Exepcted grant_type to be `urn:ietf:params:oauth:grant-type:jwt-bearer` in test")
			return
		}
		if r.Form.Get("assertion") == "" {
			t.Fatalf("Expected assertion to contain a JWT")
		}
		// TODO: validate JWT
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"scope": "openid",
			"access_token": "5301cd36-4361-4b61-98aa-0f5c3acacd21",
			"expires_in": 1799,
			"token_type": "Bearer",
			"id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJiL082T3ZWdjEreStXZ3JINVVpOVdUaW9MdDA9IiwiYWxnIjoiUlMyNTYifQ.eyJhdF9oYXNoIjoibVdETWRWRjVFNWYyRm5oYnVtT3ktQSIsInN1YiI6InRlc3RzZXJ2aWNlLnZodGVzdGFwcC52aHByb3BAdmhkZXYuaDJoZGV2b3JnLnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tIiwiYXVkaXRUcmFja2luZ0lkIjoiZWViOWIyZGItZDI0OS00NTE2LWE4NmEtMWUyMjUxYzg5Yjc0LTk3NTM3NyIsImlzcyI6Imh0dHBzOi8vZnJhdXRoYjRhNWFtLmlhbS51cy1lYXN0LnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tL29wZW5hbS9vYXV0aDIiLCJ0b2tlbk5hbWUiOiJpZF90b2tlbiIsImF1ZCI6InRlc3RzZXJ2aWNlLnZodGVzdGFwcC52aHByb3BAdmhkZXYuaDJoZGV2b3JnLnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tIiwiYXpwIjoidGVzdHNlcnZpY2Uudmh0ZXN0YXBwLnZocHJvcEB2aGRldi5oMmhkZXZvcmcucGhpbGlwcy1oZWFsdGhzdWl0ZS5jb20iLCJhdXRoX3RpbWUiOjE1MzgxMzUwMjAsInJlYWxtIjoiLyIsImV4cCI6MTUzODEzODYyMCwidG9rZW5UeXBlIjoiSldUVG9rZW4iLCJpYXQiOjE1MzgxMzUwMjB9.Jdr14sKkiOMUQRnDoceShkrE6cRLGwaSFse6lAbIEfKHp1wzDDCYu0QgL69oG_J_LbCU8ygdLmSKtww1DVt43eFdXpbKJr_n1-OarGh1aVK0lJZvx4dA2Jy_uaLpeAlt6r0ogXAO6KUTKaz_u6qZjj_DGjOO3f2WNOHqRBgfu8rqhzhViQytjPcrpFlH9YPBrZXt6j2tDfM6Ja6D8ty0E8-Qu1XUAjlO6rnvGgyjIBvAdcpVnYoeXtsG_MwAzc-oHZNANCsjmn5gpNVsU633PNpXllzPOgUEeR7z8-kT1MfZptMcRlh_L_G4FZujUTCMlSJRd4qVThWMZxR8qgtYhw"
		  }`)
	})

	err := client.ServiceLogin(*service)
	if err != nil {
		t.Fatal(err)
	}
	if !client.HasScopes("openid") {
		t.Errorf("Expected `introspect` and `cn` scope to be there")
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

func TestWithToken(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	if client.WithToken("fooz").Token() != "fooz" {
		t.Errorf("Unexpected token")
	}

}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	sharedKey := "SharedKey"
	secretKey := "SecretKey"

	client, _ = NewClient(nil, &Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      sharedKey,
		SecretKey:      secretKey,
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
		Debug:          true,
		DebugLog:       tmpfile.Name(),
	})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer client.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	client.Login("username", "password")

	fi, err := tmpfile.Stat()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("Expected something to be written to DebugLog")
	}

}
