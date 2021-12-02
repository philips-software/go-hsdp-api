package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type BucketsService struct {
	*Client
	validate *validator.Validate
}

var (
	bucketAPIVersion = "1"
)

// GetBucketOptions struct describes search criteria for looking up Bucket
type GetBucketOptions struct {
	ID            *string `url:"_id,omitempty"`
	Name          *string `url:"name,omitempty"`
	ApplicationID *string `url:"applicationId,omitempty"`
}

type Bucket struct {
	ResourceType                  string              `json:"resourceType" validate:"required"`
	ID                            string              `json:"id,omitempty"`
	Name                          string              `json:"name" validate:"required"`
	Description                   string              `json:"description"`
	PropositionID                 Reference           `json:"propositionId" validate:"required"`
	DefaultRegionID               Reference           `json:"defaultRegionId" validate:"required"`
	ReplicationRegionID           *Reference          `json:"replicationRegionId,omitempty" validate:"omitempty,dive"`
	VersioningEnabled             bool                `json:"versioningEnabled"`
	LoggingEnabled                bool                `json:"loggingEnabled"`
	CrossRegionReplicationEnabled bool                `json:"crossRegionReplicationEnabled"`
	AuditingEnabled               bool                `json:"auditingEnabled"`
	EnableCDN                     bool                `json:"enableCDN"`
	CacheControlAge               int                 `json:"cacheControlAge,omitempty" validate:"omitempty,min=300,max=1800"`
	CorsConfiguration             []CORSConfiguration `json:"corsConfiguration,omitempty" validate:"omitempty,dive"`
	Meta                          *Meta               `json:"meta,omitempty"`
}

type CORSConfiguration struct {
	AllowedOrigins []string `json:"allowedOrigins" validate:"required"`
	AllowedMethods []string `json:"allowedMethods" validate:"required,dive,oneof=GET POST PUT DELETE HEAD"`
	AllowedHeaders []string `json:"allowedHeaders,omitempty"`
	MaxAgeSeconds  int      `json:"maxAgeSeconds"`
	ExposeHeaders  []string `json:"exposeHeaders,omitempty"`
}

// Create creates a Bucket
func (c *BucketsService) Create(ac Bucket) (*Bucket, *Response, error) {
	ac.ResourceType = "Bucket"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/Bucket", ac, nil)
	req.Header.Set("api-version", bucketAPIVersion)

	var created Bucket

	resp, err := c.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

// Delete deletes the given ServiceAction
func (c *BucketsService) Delete(ac Bucket) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/Bucket/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", bucketAPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *BucketsService) GetByID(id string) (*Bucket, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/Bucket/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", bucketAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource Bucket

	resp, err := c.Do(req, &resource)
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

// Find looks up services based on GetBucketOptions
func (c *BucketsService) Find(opt *GetBucketOptions, options ...OptionFunc) (*[]Bucket, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/Bucket", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", bucketAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
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

// Update updates a standard service
func (c *BucketsService) Update(ac Bucket) (*Bucket, *Response, error) {
	ac.ResourceType = "Bucket"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/Bucket/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", bucketAPIVersion)

	var updated Bucket

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
