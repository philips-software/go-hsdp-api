package blr

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
	"net/http"
)

var (
	blobConfigurationAPIVersion = "1"
)

type BlobStorePolicy struct {
	ResourceType string                     `json:"resourceType"`
	ID           string                     `json:"id,omitempty"`
	Statement    []BlobStorePolicyStatement `json:"statement"`
}

type BlobStorePolicyStatement struct {
	Effect    string   `json:"effect"`
	Action    []string `json:"action"`
	Principal []string `json:"principal"`
	Resource  []string `json:"resource"`
}

type GetBlobStorePolicyOptions struct {
	ID *string `url:"_id,omitempty"`
}

type ConfigurationsService struct {
	*Client
	validate *validator.Validate
}

func (b *ConfigurationsService) CreateBlobStorePolicy(policy BlobStorePolicy) (*BlobStorePolicy, *Response, error) {
	policy.ResourceType = "BlobStorePolicy"
	if err := b.validate.Struct(policy); err != nil {
		return nil, nil, err
	}

	req, _ := b.NewRequest(http.MethodPost, "/configuration/BlobStorePolicy", policy, nil)
	req.Header.Set("api-version", blobConfigurationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var created BlobStorePolicy

	resp, err := b.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

func (b *ConfigurationsService) GetBlobStorePolicyByID(id string) (*Blob, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/configuration/BlobStorePolicy/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource Blob

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetByID: %w", err)
	}
	if resource.ID != id {
		return nil, nil, fmt.Errorf("returned resource does not match")
	}
	return &resource, resp, nil
}

func (b *ConfigurationsService) FindBlobStorePolicy(opt *GetBlobStorePolicyOptions, options ...OptionFunc) (*[]Blob, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/configuration/BlobStorePolicy", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobConfigurationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := b.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []Blob
	for _, c := range bundleResponse.Entry {
		var resource Blob
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

func (b *ConfigurationsService) DeleteBlobStorePolicy(policy BlobStorePolicy) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/configuration/BlobStorePolicy/"+policy.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobConfigurationAPIVersion)

	var deleteResponse interface{}

	resp, err := b.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}
