package iam

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	deviceAPIVersion = "1"
)

// CodeableConcept describes a code-able concept
type CodeableConcept struct {
	Code string `json:"code" validate:"required,min=1,max=10"`
	Text string `json:"text" validate:"max=250"`
}

// DeviceIdentifier holds device identity information
type DeviceIdentifier struct {
	System string          `json:"system" validate:"max=250"`
	Value  string          `json:"value" validate:"max=250"`
	Type   CodeableConcept `json:"type"`
}

// Device represents an IAM resource
type Device struct {
	ID                string           `json:"id,omitempty"`
	LoginID           string           `json:"loginId,omitempty" validate:"required,reserved-strings,min=5,max=50" `
	DeviceExtID       DeviceIdentifier `json:"deviceExtId" validate:"required"`
	Password          string           `json:"password,omitempty" validate:"required_without=ID,max=255"`
	Type              string           `json:"type" validate:"required,min=1,max=50"`
	RegistrationDate  *time.Time       `json:"registrationDate,omitempty"`
	ForTest           bool             `json:"forTest,omitempty"`
	IsActive          bool             `json:"isActive,omitempty"`
	DebugUntil        *time.Time       `json:"debugUntil,omitempty"`
	OrganizationID    string           `json:"organizationId" validate:"required"`
	GlobalReferenceID string           `json:"globalReferenceId" validate:"required,min=3,max=50"`
	Text              string           `json:"text,omitempty"`
	ApplicationID     string           `json:"applicationId" validate:"required"`
	Meta              *Meta            `json:"meta,omitempty"`
}

// GetDevicesOptions describes search criteria for looking up devices
type GetDevicesOptions struct {
	ID                *string `url:"_id,omitempty"`
	Count             *int    `url:"_count,omitempty"`
	Page              *int    `url:"_page,omitempty"`
	DeviceExtIDValue  *string `url:"deviceExtId.value,omitempty"`
	DeviceExtIDType   *string `url:"deviceExtId.value,omitempty"`
	DeviceExtIDSystem *string `url:"deviceExtId.system,omitempty"`
	LoginID           *string `url:"loginId,omitempty" validate:""`
	ForTest           *bool   `url:"forTest,omitempty"`
	IsActive          *bool   `url:"isActive,omitempty"`
	OrganizationID    *string `url:"organizationId,omitempty"`
	ApplicationID     *string `url:"applicationId,omitempty"`
	Type              *string `url:"type,omitempty"`
	GlobalReferenceID *string `url:"globalReferenceId,omitempty"`
	GroupID           *string `url:"groupId,omitempty"`
}

// DevicesService provides operations on IAM device resources
type DevicesService struct {
	client *Client

	validate *validator.Validate
}

// GetDevices looks up Devices based on GetDevicesOptions
// A user with DEVICE.READ permission can read device information under the user organization.
func (p *DevicesService) GetDevices(opt *GetDevicesOptions, options ...OptionFunc) (*[]Device, *Response, error) {
	req, err := p.client.newRequest(IDM, "GET", "authorize/identity/Device", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int      `json:"total"`
		Entry []Device `json:"entry"`
	}

	resp, err := p.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &bundleResponse.Entry, resp, err
}

// GetDeviceByID retrieves a device by ID
func (p *DevicesService) GetDeviceByID(deviceID string) (*Device, *Response, error) {
	devices, resp, err := p.GetDevices(&GetDevicesOptions{
		ID: &deviceID,
	})
	if devices == nil || len(*devices) == 0 {
		return nil, resp, fmt.Errorf("GetDeviceByID: %v %w", err, ErrNotFound)
	}
	return &(*devices)[0], resp, err
}

// CreateDevice creates a Device
// A user with DEVICE.WRITE permission can create devices under the organization.
func (p *DevicesService) CreateDevice(device Device) (*Device, *Response, error) {
	if err := p.validate.Struct(device); err != nil {
		return nil, nil, err
	}
	req, _ := p.client.newRequest(IDM, "POST", "authorize/identity/Device", device, nil)
	req.Header.Set("api-version", deviceAPIVersion)

	var createdDevice Device

	resp, err := p.client.do(req, &createdDevice)
	if resp == nil {
		return nil, nil, ErrOperationFailed
	}
	ok := resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated
	if !ok {
		return nil, resp, err
	}
	var id string
	count, _ := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/Device/%s", &id)
	if count == 0 {
		return nil, resp, ErrCouldNoReadResourceAfterCreate
	}
	return p.GetDeviceByID(id)
}

// UpdateDevice updates Device properties.
// Any user with DEVICE.WRITE permission within the organization can update device properties.
// The entire resource data must be passed as request body to update a device.
// If read-only attributes (such as id, loginId, password, meta, organizationId) are passed, that will be ignored.
func (p *DevicesService) UpdateDevice(device Device) (*Device, *Response, error) {
	req, err := p.client.newRequest(IDM, "PUT", "authorize/identity/Device/"+device.ID, &device, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", deviceAPIVersion)

	var updatedDevice Device

	resp, err := p.client.do(req, &updatedDevice)
	if err != nil {
		return nil, resp, err
	}
	return &updatedDevice, resp, err

}

// DeleteDevice deletes the given Device
// The is usually done by a organization administrator.
// Any user with DEVICE.WRITE or DEVICE.DELETE permission within
// the organization can delete a device from an organization.
func (p *DevicesService) DeleteDevice(device Device) (bool, *Response, error) {
	req, err := p.client.newRequest(IDM, "DELETE", "authorize/identity/Device/"+device.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", deviceAPIVersion)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

// ChangePassword changes the password. The current pasword must be provided as well.
// No password history will be maintained for device.
func (p *DevicesService) ChangePassword(deviceID, oldPassword, newPassword string) (bool, *Response, error) {
	body := struct {
		OldPassword string `json:"oldPassword" validate:"required,min=8"`
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	if err := p.validate.Struct(body); err != nil {
		return false, nil, err
	}
	return p.deviceActionV(deviceID, body, "$change-password", deviceAPIVersion)
}

func (p *DevicesService) deviceActionV(deviceID string, body interface{}, action, apiVersion string) (bool, *Response, error) {
	req, err := p.client.newRequest(IDM, "POST", "authorize/identity/Device/"+deviceID+"/"+action, body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", apiVersion)

	var bundleResponse interface{}

	resp, err := p.client.doSigned(req, &bundleResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}
