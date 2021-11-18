package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type FirmwareComponentVersionsService struct {
	*Client
	validate *validator.Validate
}

var (
	firmwareComponentVersionAPIVersion = "1"
)

type Fingerprint struct {
	Algorithm string `json:"alg" validate:"required"`
	Hash      string `json:"hash" validate:"required"`
}

type EncryptionInfo struct {
	Encrypted     bool   `json:"encrypted" validate:"required"`
	Algorithm     string `json:"alg,omitempty"`
	DecryptionKey string `json:"decryptionKey,omitempty"`
}

type FirmwareComponentVersion struct {
	ResourceType               string          `json:"resourceType" validate:"required"`
	ID                         string          `json:"id,omitempty"`
	Version                    string          `json:"version" validate:"required,min=1,max=20"`
	Description                string          `json:"description" validate:"omitempty,max=250"`
	Size                       int             `json:"size"`
	BlobURL                    string          `json:"blobUrl" validate:"omitempty,max=1024"`
	ComponentRequired          bool            `json:"componentRequired"`
	FingerPrint                *Fingerprint    `json:"fingerPrint,omitempty" validate:"omitempty,dive"`
	EncryptionInfo             *EncryptionInfo `json:"encryptionInfo,omitempty" validate:"omitempty,dive"`
	FirmwareComponentId        Reference       `json:"firmwareComponentId" validate:"required,dive"`
	PreviousComponentVersionId *Reference      `json:"previousComponentVersionId,omitempty" validate:"omitempty,dive"`
	EffectiveDate              string          `json:"effectiveDate" validate:"required"`
	DeprecatedDate             string          `json:"deprecatedDate,omitempty"`
	CustomResource             json.RawMessage `json:"customResource,omitempty" validate:"omitempty,max=2048"`
	Meta                       *Meta           `json:"meta,omitempty"`
}

// GetFirmwareComponentVersionOptions struct describes search criteria for looking up FirmwareComponentVersion
type GetFirmwareComponentVersionOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

// Create creates a FirmwareComponentVersion
func (c *FirmwareComponentVersionsService) Create(ac FirmwareComponentVersion) (*FirmwareComponentVersion, *Response, error) {
	ac.ResourceType = "FirmwareComponentVersion"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/FirmwareComponentVersion", ac, nil)
	req.Header.Set("api-version", firmwareComponentVersionAPIVersion)

	var created FirmwareComponentVersion

	resp, err := c.Do(req, &created)

	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	if !ok {
		return nil, resp, err
	}
	if resp == nil {
		return nil, resp, fmt.Errorf("create (resp=nil): %w", ErrCouldNoReadResourceAfterCreate)
	}

	return c.GetByID(created.ID)
}

// Delete deletes the given FirmwareComponentVersion
func (c *FirmwareComponentVersionsService) Delete(ac FirmwareComponentVersion) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/FirmwareComponentVersion/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", firmwareComponentVersionAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *FirmwareComponentVersionsService) GetByID(id string) (*FirmwareComponentVersion, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/FirmwareComponentVersion/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentVersionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource FirmwareComponentVersion

	resp, err := c.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetByID: %w", err)
	}
	return &resource, resp, nil
}

// Find looks up services based on GetFirmwareComponentVersionOptions
func (c *FirmwareComponentVersionsService) Find(opt *GetFirmwareComponentVersionOptions, options ...OptionFunc) (*[]FirmwareComponentVersion, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/FirmwareComponentVersion", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentVersionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []FirmwareComponentVersion
	for _, c := range bundleResponse.Entry {
		var resource FirmwareComponentVersion
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *FirmwareComponentVersionsService) Update(ac FirmwareComponentVersion) (*FirmwareComponentVersion, *Response, error) {
	ac.ResourceType = "FirmwareComponentVersion"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/FirmwareComponentVersion/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", firmwareComponentVersionAPIVersion)

	var updated FirmwareComponentVersion

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
