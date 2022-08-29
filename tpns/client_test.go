package tpns

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	muxTPNS    *http.ServeMux
	serverTPNS *httptest.Server
	tpnsClient *Client
)

func setup(t *testing.T) func() {
	muxTPNS = http.NewServeMux()
	serverTPNS = httptest.NewServer(muxTPNS)

	tpnsClient, _ = NewClient(nil, &Config{
		TPNSURL: serverTPNS.URL,
	})

	muxTPNS.HandleFunc("/tpns/PushMessage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
		}
		if apiVersion := r.Header.Get("Api-Version"); apiVersion != "2" {
			t.Errorf("Expected Api-Version = 2, got %s", apiVersion)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var tpnsRequest Message
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &tpnsRequest)
		if err != nil {
			t.Errorf("Invalid body in request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tpnsRequest.MessageType == "" {
			t.Errorf("Empty MessageType")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"issue": [
			  {
				"Severity": "information",
				"Code": {
				  "coding": [
					{
					  "system": "MS",
					  "code": "201"
					}
				  ]
				},
				"Details": "Notification Sent"
			  }
			]
		  }`)
	})

	return func() {
		serverTPNS.Close()
		tpnsClient.Close()
	}
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	tpnsClient, _ = NewClient(nil, &Config{
		TPNSURL:  serverTPNS.URL,
		Debug:    true,
		DebugLog: tmpfile.Name(),
	})
	defer func() {
		tpnsClient.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	_, _, _ = tpnsClient.Messages.Push(&Message{
		PropositionID: "XYZ",
		MessageType:   "Push",
		Content:       "YAY!",
		Targets:       []string{"foo"},
	})
	fi, err := tmpfile.Stat()
	if !assert.Nil(t, err) {
		return
	}
	assert.Less(t, int64(0), fi.Size())
}

func TestPush(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ok, resp, err := tpnsClient.Messages.Push(&Message{
		PropositionID: "XYZ",
		MessageType:   "Push",
		Content:       "YAY!",
		Targets:       []string{"foo"},
	})

	if !ok {
		t.Errorf("Expected call to succeed: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK")
		return
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
