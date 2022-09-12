package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestDevicesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	managingOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9a"
	deviceID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc1"
	muxIDM.HandleFunc("/authorize/identity/Device", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Location", "/authorize/identity/Device/"+deviceID)
			w.WriteHeader(http.StatusCreated)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "total": 1,
  "entry": [
    {
      "loginId": "andydevice",
      "organizationId": "`+managingOrgID+`",
      "applicationId": "711171ab-d28c-4616-a314-f95584e280c3",
      "deviceExtId": {
        "type": {
          "code": "ID"
        },
        "system": "http://www.philips.co.id/c-m-ho/fake/fakedevice",
        "value": "0001"
      },
      "type": "Device",
      "isActive": true,
      "globalReferenceId": "c157bd2e-e992-4b5e-88ab-911766b7b8f4",
      "id": "`+deviceID+`",
      "meta": {
        "versionId": "0",
        "lastModified": "2020-08-23T20:47:374.040Z"
      }
    }
  ]
}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Device/"+deviceID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
      "loginId": "andydevice",
      "organizationId": "`+managingOrgID+`",
      "applicationId": "711171ab-d28c-4616-a314-f95584e280c3",
      "deviceExtId": {
        "type": {
          "code": "ID"
        },
        "system": "http://www.philips.co.id/c-m-ho/fake/fakedevice",
        "value": "0001"
      },
      "type": "Device",
      "isActive": true,
      "globalReferenceId": "c157bd2e-e992-4b5e-88ab-911766b7b8f4",
      "id": "`+deviceID+`",
      "meta": {
        "versionId": "1",
        "lastModified": "2020-08-23T20:47:374.040Z"
      }
    }`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Device/"+deviceID+"/$change-password", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	d := Device{
		LoginID:           "andydevice",
		Password:          "SecretPasw0rd!",
		GlobalReferenceID: "c157bd2e-e992-4b5e-88ab-911766b7b8f4",
		ApplicationID:     "711171ab-d28c-4616-a314-f95584e280c3",
		OrganizationID:    managingOrgID,
		Type:              "Device",
		IsActive:          true,
		DeviceExtID: DeviceIdentifier{
			Type: CodeableConcept{
				Code: "ID",
				Text: "Device identification",
			},
			System: "http://www.philips.co.id/c-m-ho/fake/fakedevice",
			Value:  "0001",
		},
	}

	device, resp, err := client.Devices.CreateDevice(d)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, device) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode()) {
		return
	}
	assert.Equal(t, managingOrgID, device.OrganizationID)

	device, resp, err = client.Devices.UpdateDevice(*device)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, device) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode()) {
		return
	}
	assert.Equal(t, managingOrgID, device.OrganizationID)

	device, resp, err = client.Devices.GetDeviceByID(deviceID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, device) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode()) {
		return
	}
	assert.Equal(t, deviceID, device.ID)
	ok, resp, err := client.Devices.DeleteDevice(*device)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusNoContent, resp.StatusCode()) {
		return
	}
	assert.Equal(t, true, ok)

	ok, resp, err = client.Devices.ChangePassword(device.ID, "foofoo12", "barbar12")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.Equal(t, http.StatusNoContent, resp.StatusCode()) {
		return
	}
	assert.Equal(t, true, ok)
}

func TestValidationDevices(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	_, _, err := client.Devices.CreateDevice(Device{
		LoginID: "withdelegations",
	})
	if !assert.NotNil(t, err) {
		return
	}
	assert.IsType(t, validator.ValidationErrors{}, err)

	_, _, err = client.Devices.CreateDevice(Device{})
	if !assert.NotNil(t, err) {
		return
	}

	_, _, err = client.Devices.ChangePassword("id", "foo", "tooshort")
	if !assert.NotNil(t, err) {
		return
	}
}
