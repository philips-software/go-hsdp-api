package stl_test

import (
	"context"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
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
	v.HSDPProductKey = "key"
	ok, err = v.Validate()
	assert.NotNil(t, err)
	assert.False(t, ok)
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
