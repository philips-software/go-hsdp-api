package logging

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hsdp/go-signer"
	"github.com/jeffail/gabs"
)

var (
	muxLogger     *http.ServeMux
	serverLogger  *httptest.Server
	client        *Client
	validResource = Resource{
		ID:                  "deb545e2-ccea-4868-99fe-b9dfbf5ce56e",
		ResourceType:        "LogEvent",
		ServerName:          "foo.bar.com",
		ApplicationName:     "some-space",
		EventID:             "1",
		Category:            "Tracelog",
		Component:           "PHS",
		TransactionID:       "5bc4ce05-37b5-4f08-89e4-ed73790f8058",
		ServiceName:         "mcvs",
		ApplicationInstance: "85e597cb-2648-4187-78ec-2c58",
		ApplicationVersion:  "0.0.0",
		OriginatingUser:     "ActiveUser",
		LogTime:             "2017-10-15T01:53:20Z",
		Severity:            "INFO",
		LogData: LogData{
			Message: "hello world",
		},
	}
)

const (
	productKey   = "859722b3-64dd-4be8-a522-c2dbf88c86b5"
	sharedKey    = "SharedKey"
	sharedSecret = "SharedSecret"
)

func setup(t *testing.T, config Config) func() {
	muxLogger = http.NewServeMux()
	serverLogger = httptest.NewServer(muxLogger)

	config.BaseURL = serverLogger.URL

	client, _ = NewClient(nil, config)

	muxLogger.HandleFunc("/core/log/LogEvent", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s, _ := signer.New(sharedKey, sharedSecret)
		if ok, _ := s.ValidateRequest(r); !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		j, _ := gabs.ParseJSON(body)
		pk, ok := j.Path("productKey").Data().(string)
		if !ok {
			t.Errorf("Missing productKey field")
		}
		if pk != productKey {
			w.WriteHeader(http.StatusUnprocessableEntity)
			io.WriteString(w, `{
				"issue": [
					{
						"severity": "error",
						"code": "value",
						"details": {
							"coding": [
								{
									"system": "https://www.hl7.org/fhir/valueset-operation-outcome.html",
									"code": "MSG_PARAM_INVALID"
								}
							],
							"text": "Invalid parameter value"
						},
						"diagnostics": "Invalid parameter value"
					}
				],
				"resourceType": "OperationOutcome"
			}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	return func() {
		serverLogger.Close()
	}
}

func TestStoreResources(t *testing.T) {
	teardown := setup(t, Config{
		SharedKey:    sharedKey,
		SharedSecret: sharedSecret,
		ProductKey:   productKey,
	})
	defer teardown()

	var resource = []Resource{
		validResource,
	}

	resp, err := client.StoreResources(resource, len(resource))
	if err != nil {
		t.Errorf("Unexpected response: %v", err)
		return
	}
	if resp == nil {
		t.Errorf("Unexpected nil value for response")
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP 201, Got: %d", resp.StatusCode)
	}
}

func TestStoreResourcesWithInvalidKey(t *testing.T) {
	teardown := setup(t, Config{
		SharedKey:    sharedKey,
		SharedSecret: sharedSecret,
		ProductKey:   "089db3e5-3e3e-4445-8903-29cc848194b1",
	})
	defer teardown()

	var resource = []Resource{
		validResource,
	}

	resp, err := client.StoreResources(resource, len(resource))
	if err == nil {
		t.Errorf("Expected error message")
	}
	if resp == nil {
		t.Errorf("Unexpected nil value for response")
		return
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected HTTP %d, Got: %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestStoreResourcesWithInvalidKeypair(t *testing.T) {
	os.Setenv("DEBUG", "true")
	teardown := setup(t, Config{
		SharedKey:    "bogus",
		SharedSecret: "keys",
		ProductKey:   productKey,
	})
	defer teardown()

	var resource = []Resource{
		validResource,
	}

	resp, err := client.StoreResources(resource, len(resource))
	if resp == nil {
		t.Errorf("Unexpected nil value for response")
		return
	}
	_ = err.Error() // Just to up coverage
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected HTTP 403, Got: %d", resp.StatusCode)
	}
	if err == nil {
		t.Errorf("Expected error response")
	}
}
