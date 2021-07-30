package stl_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAppLoggingInputValidate(t *testing.T) {
	v := stl.UpdateAppLoggingInput{}
	ok, err := v.Validate()
	assert.Nil(t, err)
	assert.True(t, ok)
	v.RawConfig = "[OUTPUT]"
	ok, err = v.Validate()
	assert.Nil(t, err)
	assert.True(t, ok)
	v = stl.UpdateAppLoggingInput{
		SerialNumber: "xxx",
		AppLogging: stl.AppLogging{
			HSDPProductKey:   "a",
			HSDPSecretKey:    "b",
			HSDPIngestorHost: "http://foo",
			HSDPSharedKey:    "c",
			HSDPCustomField:  nil,
		},
	}
	ok, err = v.Validate()
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestUpdateAppFirewallExceptions(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	serial := "A444900Z0822111"

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "updateAppFirewallException": {
      "statusCode": 200,
      "success": true,
      "message": "Successfully updated app firewall exceptions",
      "appFirewallException": {
        "deviceId": 53615,
        "tcp": [
          8080,
          443
        ],
        "udp": []
      }
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	fw, err := client.Config.UpdateAppFirewallExceptions(ctx, stl.UpdateAppFirewallExceptionInput{
		SerialNumber: serial,
		AppFirewallException: stl.AppFirewallException{
			TCP: []int{8080, 443},
			UDP: []int{},
		},
	})
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, len(fw.TCP))
	assert.Equal(t, 0, len(fw.UDP))
	assert.Equal(t, 8080, fw.TCP[0])
	assert.Equal(t, int64(53615), fw.DeviceID)
}

func TestGetFirewallExceptionsBySerial(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	serial := "A444900Z0822111"

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "appFirewallException": {
      "deviceId": 53615,
      "tcp": [
        4443,
        80,
        8080,
        443
      ],
      "udp": []
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	fw, err := client.Config.GetFirewallExceptionsBySerial(ctx, serial)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, fw) {
		return
	}
	assert.Equal(t, 4, len(fw.TCP))
	assert.Equal(t, 0, len(fw.UDP))
	assert.Equal(t, 4443, fw.TCP[0])
	assert.Equal(t, int64(53615), fw.DeviceID)
}

func TestGetAppLoggingBySerial(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	serial := "A444900Z0822111"

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "appLogging": {
      "deviceId": 53615,
      "rawConfig": "[OUTPUT]",
      "hsdpLogging": false,
      "hsdpIngestorHost": "",
      "hsdpSharedKey": "",
      "hsdpSecretKey": "",
      "hsdpProductKey": "",
      "hsdpCustomField": false
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	logging, err := client.Config.GetAppLoggingBySerial(ctx, serial)
	assert.Nil(t, err)
	if !assert.NotNil(t, logging) {
		return
	}
	assert.Equal(t, "[OUTPUT]", logging.RawConfig)
}

func TestUpdateAppLogging(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()
	serial := "A444900Z0822111"

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "updateAppLogging": {
      "statusCode": 200,
      "success": true,
      "message": "Successfully updated app logging config",
      "appLogging": {
        "deviceId": 53615,
        "rawConfig": "[OUTPUT]",
        "hsdpLogging": false,
        "hsdpIngestorHost": "",
        "hsdpSharedKey": "",
        "hsdpSecretKey": "",
        "hsdpProductKey": "",
        "hsdpCustomField": false
      }
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	logging, err := client.Config.UpdateAppLogging(ctx, stl.UpdateAppLoggingInput{
		SerialNumber: serial,
		AppLogging: stl.AppLogging{
			RawConfig: "[OUTPUT]",
		},
	})
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "[OUTPUT]", logging.RawConfig)
}

func TestConfigServiceFailures(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	serial := "A444900Z0822111"
	ctx := context.Background()
	l, err := client.Config.UpdateAppLogging(ctx, stl.UpdateAppLoggingInput{})
	assert.NotNil(t, err)
	assert.Nil(t, l)
	l, err = client.Config.GetAppLoggingBySerial(ctx, serial)
	assert.NotNil(t, err)
	assert.Nil(t, l)
	f, err := client.Config.UpdateAppFirewallExceptions(ctx, stl.UpdateAppFirewallExceptionInput{})
	assert.NotNil(t, err)
	assert.Nil(t, f)
	f, err = client.Config.GetFirewallExceptionsBySerial(ctx, serial)
	assert.NotNil(t, err)
	assert.Nil(t, f)
}
