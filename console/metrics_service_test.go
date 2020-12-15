package console_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/console"

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
        "guid": "c3971808-c6e2-487d-9bb2-20c116ad03a7",
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

	muxCONSOLE.HandleFunc("/v3/metrics/instances/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
        "createdAt": "2018-08-24T23:32:05Z",
        "guid": "c3971808-c6e2-487d-9bb2-20c116ad03a7",
        "name": "metrics",
        "organization": "system",
        "space": "service-brokers"
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
	_, _, err = client.Metrics.GetInstanceByID("2")
	assert.NotNil(t, err)
	instance, resp, err := client.Metrics.GetInstanceByID("1")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.NotNil(t, instance)
}

func TestMetricsRuleCalls(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxCONSOLE.HandleFunc("/v3/metrics/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "status": "success",
  "data": {
    "groups": [
      {
        "annotations": {
          "description": "description",
          "resolved": "",
          "summary": "summary"
        },
        "description": "RabbitMQ instance load is high",
        "id": "rabbit_load_is_high",
        "metric": "node_load1",
        "rule": {
          "extras": null,
          "operators": [
            ">",
            ">="
          ],
          "subject": "RabbitMQ instance",
          "threshold": {
            "default": 75,
            "max": 100,
            "min": 0,
            "type": "range",
            "unit": [
              "%"
            ]
          }
        },
        "template": "template"
      }
    ]
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCONSOLE.HandleFunc("/v3/metrics/rules/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
        "annotations": {
          "description": "description",
          "resolved": "",
          "summary": "summary"
        },
        "description": "RabbitMQ instance load is high",
        "id": "rabbit_load_is_high",
        "metric": "node_load1",
        "rule": {
          "extras": null,
          "operators": [
            ">",
            ">="
          ],
          "subject": "RabbitMQ instance",
          "threshold": {
            "default": 75,
            "max": 100,
            "min": 0,
            "type": "range",
            "unit": [
              "%"
            ]
          }
        },
        "template": "template"
      }`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	groups, resp, err := client.Metrics.GetGroupedRules()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.Equal(t, 1, len(*groups))
	_, _, err = client.Metrics.GetRuleByID("2")
	assert.NotNil(t, err)
	rule, resp, err := client.Metrics.GetRuleByID("1")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.NotNil(t, rule)
}

func TestAutoscalerCalls(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	instanceID := "c20c76b7-6622-4fc0-892b-92c0caff91a5"

	muxCONSOLE.HandleFunc("/v3/metrics/"+instanceID+"/autoscalers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "status": "success",
  "data": {
    "application": {
      "name": "consul-app",
      "enabled": true,
      "minInstances": 1,
      "maxInstances": 10,
      "thresholds": [
        {
          "name": "cpu",
          "enabled": true,
          "min": 5,
          "max": 95
        },
        {
          "name": "memory",
          "enabled": false,
          "min": 20,
          "max": 100
        },
        {
          "name": "http-latency",
          "enabled": false,
          "min": 0.01,
          "max": 10
        },
        {
          "name": "http-rate",
          "enabled": false,
          "min": 5,
          "max": 1000000
        }
      ]
    }
  }
}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "status": "success",
  "data": {
    "applications": [
      {
        "name": "consul-app",
        "enabled": true,
        "minInstances": 1,
        "maxInstances": 10
      },
      {
        "name": "mds-app",
        "enabled": false,
        "minInstances": 1,
        "maxInstances": 10
      }
    ]
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCONSOLE.HandleFunc("/v3/metrics/"+instanceID+"/autoscalers/consul-app", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "status": "success",
  "data": {
    "application": {
      "name": "consul-app",
      "enabled": true,
      "minInstances": 1,
      "maxInstances": 10,
      "thresholds": [
        {
          "name": "cpu",
          "enabled": true,
          "min": 5,
          "max": 95
        },
        {
          "name": "memory",
          "enabled": false,
          "min": 20,
          "max": 100
        },
        {
          "name": "http-latency",
          "enabled": false,
          "min": 0.01,
          "max": 10
        },
        {
          "name": "http-rate",
          "enabled": false,
          "min": 5,
          "max": 1000000
        }
      ]
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	_, _, err = client.Metrics.GetApplicationAutoscalers("xx")
	assert.NotNil(t, err)
	apps, resp, err := client.Metrics.GetApplicationAutoscalers(instanceID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	assert.Equal(t, 2, len(*apps))
	app, resp, err := client.Metrics.GetApplicationAutoscaler(instanceID, "consul-app")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	if !assert.NotNil(t, app) {
		return
	}
	assert.Equal(t, "consul-app", app.Name)

	_, _, err = client.Metrics.UpdateApplicationAutoscaler("xx", *app)
	assert.NotNil(t, err)
	app, resp, err = client.Metrics.UpdateApplicationAutoscaler(instanceID, *app)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}
	if !assert.NotNil(t, app) {
		return
	}
	assert.Equal(t, "consul-app", app.Name)
}

func TestRetryableError(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	instanceID := "c20c76b7-6622-4fc0-892b-92c0caff91a5"

	muxCONSOLE.HandleFunc("/v3/metrics/"+instanceID+"/autoscalers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{
  "status": "error",
  "error": {
    "code": "BAD_REQUEST",
    "message": "invalid character u003c looking for beginning of value"
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
	app, resp, err := client.Metrics.UpdateApplicationAutoscaler(instanceID, console.Application{
		Enabled:      true,
		MinInstances: 1,
		MaxInstances: 10,
		Name:         "foo",
	})
	assert.NotNil(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Nil(t, app) {
		return
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	assert.Contains(t, resp.Error.Message, "invalid character")
}
