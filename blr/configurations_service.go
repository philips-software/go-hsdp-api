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

type Bucket struct {
	ResourceType                 string            `json:"resourceType" validate:"required"`
	ID                           string            `json:"id,omitempty"`
	Name                         string            `json:"name" validate:"required"`
	EnableHSDPDomain             bool              `json:"enableHSDPDomain"`
	EnableCDN                    bool              `json:"enableCDN"`
	PriceClass                   string            `json:"priceClass,omitempty"`
	CacheControlAge              int               `json:"cacheControlAge,omitempty"`
	PropositionID                Reference         `json:"propositionId" validate:"required"`
	CorsConfiguration            CorsConfiguration `json:"corsConfiguration"`
	EnableCreateOrDeleteBlobMeta bool              `json:"enableCreateOrDeleteBlobMeta"`
}

type Reference struct {
	Reference string `json:"reference"`
	Display   string `json:"display"`
}

// GetBucketOptions struct describes search criteria for looking up Bucket
type GetBucketOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	PropositionID *string `url:"propositionId,omitempty"`
	Count         *int    `url:"_count,omitempty"`
	Page          *int    `url:"page,omitempty"`
}

type CorsConfiguration struct {
	AllowedOrigins []string `json:"allowedOrigins"`
	AllowedMethods []string `json:"allowedMethods"`
	AllowedHeaders []string `json:"allowedHeaders"`
	ExposeHeaders  []string `json:"exposeHeaders"`
	MaxAgeSeconds  int      `json:"maxAgeSeconds"`
}

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

func (b *ConfigurationsService) GetBlobStorePolicyByID(id string) (*BlobStorePolicy, *Response, error) {
	policies, resp, err := b.FindBlobStorePolicy(&GetBlobStorePolicyOptions{
		ID: &id,
	})
	if err != nil {
		return nil, nil, err
	}
	if len(*policies) == 0 {
		return nil, nil, fmt.Errorf("policy with id '%s' not found", id)
	}
	return &(*policies)[0], resp, nil
}

func (b *ConfigurationsService) FindBlobStorePolicy(opt *GetBlobStorePolicyOptions, options ...OptionFunc) (*[]BlobStorePolicy, *Response, error) {
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
	var resources []BlobStorePolicy
	for _, c := range bundleResponse.Entry {
		var resource BlobStorePolicy
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

// Buckets

func (b *ConfigurationsService) CreateBucket(bucket Bucket) (*Bucket, *Response, error) {
	bucket.ResourceType = "Bucket"
	if err := b.validate.Struct(bucket); err != nil {
		return nil, nil, err
	}

	req, _ := b.NewRequest(http.MethodPost, "/configuration/Bucket", bucket, nil)
	req.Header.Set("api-version", blobConfigurationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var created Bucket

	resp, err := b.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

// UpdateBucket updates a bucket
func (b *ConfigurationsService) UpdateBucket(bucket Bucket) (*Bucket, *Response, error) {
	bucket.ResourceType = "Bucket"
	if err := b.validate.Struct(bucket); err != nil {
		return nil, nil, err
	}
	req, _ := b.NewRequest(http.MethodPut, "/configuration/Bucket/"+bucket.ID, bucket, nil)
	req.Header.Set("api-version", blobConfigurationAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var updated Bucket

	resp, err := b.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}

func (b *ConfigurationsService) DeleteBucket(bucket Bucket) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/configuration/Bucket/"+bucket.ID, nil, nil)
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

func (b *ConfigurationsService) GetBucketByID(id string) (*Bucket, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/configuration/Bucket/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource Bucket

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

func (b *ConfigurationsService) FindBucket(opt *GetBucketOptions, options ...OptionFunc) (*[]Bucket, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/configuration/Bucket", opt, options...)
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
	var resources []Bucket
	for _, c := range bundleResponse.Entry {
		var resource Bucket
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}
