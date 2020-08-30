package console

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsInstanceCalls(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxCONSOLE.HandleFunc("/v3/metrics/instances", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "instances": [
      {
        "createdAt": "2018-08-24T23:32:05Z",
        "guid": "670f1fa9-f40e-449c-be77-f383d72cc7f6",
        "name": "metrics",
        "organization": "system",
        "space": "service-brokers"
      }
    ]
  },
  "status": "success"
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	instances, resp, err := client.Metrics.GetInstances()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.Equal(t, 1, len(*instances))
}
