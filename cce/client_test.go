package cce_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/cce"

	"github.com/stretchr/testify/assert"
)

var (
	muxIAM    *http.ServeMux
	serverIAM *httptest.Server
	muxCCE    *http.ServeMux
	serverCCE *httptest.Server

	cceClient *cce.Client

	privateKey = "-----BEGIN RSA PRIVATE KEY-----MIIEpAIBAAKCAQEAwM8LhQS4OB6e0xrMHE20NI/vWAwdgG3eoa50mlhlDwKQg0/sMYUKZBHkcit4rEQvgpXb36WtBhLAGC5gxLCBioRMfFG6c+DS9xyKXCexTTQZC1qBZlh1M7kq6oywnqfozBJ/9nAneOIkqA4NT9sy7jSMDuGFursL7p0iB1LrqEptBxm1zZKOw9GXUzqGTa+jdVj4DoviBtm6DCnQ61ucOEkl6DGvll5QBI693XIomqIbBICRHeMcTNoJ2GmKPYRITazKyk7FJc7Sn7E5T+ZB9StkX4BgiZjjZwnUKYpYGX65tipy3nrTzzZHfM83Rl91Gn1fbWIGxipxhp544XcjpwIDAQABAoIBAQC81YrYumiaPhMbenFRfyDxIc8uEp+KOxEClNQKnmxLqR1UHiCb10r3+zYcQ0sqnJVTdeYkQiUVf6O3iySnPp+AxFYMpBbSiuzTrJYt74oMrOuiXP/C9vvCrqXDlgsdOCIeTDgbanieQg3YsfqDrZFSDxDlOic5XRwwlKDRP3siFBuLrZZ/PtzNO4uMeBSKxGveghdQJiQ07XZoZ7bBsRU5lIV7tI+bpryA+xLu6C/LSRgtnwvHxnrppdvVOB4ZMc1IcigmCvMKCdUG1AQuxFH/v8ACuMuDubYRCJECtKlvQHNryA/uF8FYINjWDFiyoUmg2uu+Xtk13dY2uJoEycshAoGBAP1C1xs/9gaji6kwMP449IPH0pS24Nm0Rv3nxN6BPp0N+WpqCs9OAIFeNUZzScI6rMUtOBxEbddLGCkXCWnpvaPUoApx0192aW7mZ3mfhYQsaHohQHjkrzOCuRULaEbyET//w+7FvmxRHrjrEkuu8Tg06Uxt9394f5qVCrx6Eo2ZAoGBAMLk10htsuHx1KwYw7r6SLGaHhV60viUfKhiVXWUZiYCfgBtcr8Q/XBNFy/7e6Sm5o5xDBknFFcraMFKXpr9UMyOLLZLRqh8rxdRtoI3bGR1pNnx1BjkF6JmbDn8zqcibTet8nRXf7HoozolpbqL0QF8IXPzZpLd7o/+4RF4mYM/AoGBAOxtgpxoyIeIE/A9Ee+yQenoGFlGpH/4QTH1NR827rn1erryBedjfStIRFnhdKEC35kvTqts4lHTQ9nQLLSYRbZ033cArf/3bhPeugibeCxcvKgO9L4nVruytI/F13IrtxjU7xevuMYrsI+Wu7y1s3DyTD1Sh3OTjSRFMQGkwD85AoGANrfJOayS7JzY+Ph6+6QJhNOgXqd9VA1ccmopVDm19DX+6l/QN5StkzoRqIcSz8eMM7HJk8ZFD7RAVQRsS1eTt9qy8vtvex6GiiWG+EhXRl1BS295/QMNH6th92XjH0mrIFbWG5P1Zh3KtiibvyRCKgiP294ajmGA+Sy2RBF4CEECgYAKlx5gd2/1iCjOtnWRMjzSfhPUKe6W7FE3QwgoT7ZTA2xM4gLsSfC148G2EZMjDdKrHUdLC+hAQ4tRAqLMLjfyIhU/SKMM4NnVZlcZqdbWJP2BviXUz1rGKlwyt8fhvzbeiRbqTQxDdl0cjpE4Zi8aFSyczOGBsxQ0X57n5ARAtA==-----END RSA PRIVATE KEY-----"
	serviceID  = "testservice.vhtestapp.vhprop@vhdev.h2hdevorg.philips-healthsuite.com"
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxCCE = http.NewServeMux()
	serverCCE = httptest.NewServer(muxCCE)

	var err error

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxCCE.HandleFunc("/"+cce.DiscoveryPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := cce.DiscoveryEndpoints{
			TokenEndpoint:         serverIAM.URL + "/authorize/oauth2/token",
			IntrospectionEndpoint: serverCCE.URL + "/introspect",
			DiscoveryEndpoint:     serverCCE.URL + "/discovery",
		}
		data, _ := json.Marshal(response)
		_, _ = io.WriteString(w, string(data))

	})
	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
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

	cceClient, err = cce.NewClient(nil, &cce.Config{
		ServiceID:  serviceID,
		PrivateKey: privateKey,
		BaseURL:    serverCCE.URL,
	})
	assert.Nilf(t, err, "failed to create cceClient: %v", err)

	return func() {
		serverIAM.Close()
		serverCCE.Close()
	}
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	cceClient, err = cce.NewClient(nil, &cce.Config{
		ServiceID:  serviceID,
		PrivateKey: privateKey,
		BaseURL:    serverCCE.URL,
		Debug:      true,
		DebugLog:   tmpfile.Name(),
	})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer cceClient.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	_, _, err = cceClient.Discovery()
	assert.NotNil(t, err)

	fi, err := tmpfile.Stat()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, fi.Size(), "Expected something to be written to DebugLog")
}
