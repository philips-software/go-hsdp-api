package stl_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetDevices(t *testing.T) {
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
    "device": {
      "id": 53615,
      "name": "Andy SME100-1",
      "state": "authorized",
      "region": "na1",
      "serialNumber": "`+serial+`",
      "primaryInterface": {
        "name": "br0",
        "address": "192.168.2.2"
      }
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	device, err := client.Devices.GetDeviceBySerial(ctx, serial)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, device) {
		return
	}
	assert.Equal(t, serial, device.SerialNumber)
	assert.Equal(t, int64(53615), device.ID)
	assert.Equal(t, "192.168.2.2", device.PrimaryInterface.Address)
	assert.Equal(t, "Andy SME100-1", device.Name)

	device, err = client.Devices.GetDeviceByID(ctx, 53615)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, device) {
		return
	}
	assert.Equal(t, "Andy SME100-1", device.Name)
}

func TestSyncDevice(t *testing.T) {
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
    "syncDeviceConfigs": {
      "statusCode": 200,
      "success": true,
      "message": "Successfully sent command to synchronize configs."
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	err = client.Devices.SyncDeviceConfig(ctx, serial)
	if !assert.Nil(t, err) {
		return
	}
}

func TestDeviceFailures(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	serial := "A444900Z0822111"
	ctx := context.Background()

	device, err := client.Devices.GetDeviceByID(ctx, 1)
	assert.NotNil(t, err)
	assert.Nil(t, device)
	device, err = client.Devices.GetDeviceByID(ctx, 1)
	assert.NotNil(t, err)
	assert.Nil(t, device)
	err = client.Devices.SyncDeviceConfig(ctx, serial)
	assert.NotNil(t, err)
}
