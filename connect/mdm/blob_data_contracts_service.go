package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type BlobDataContractsService struct {
	*Client
	validate *validator.Validate
}

var (
	blobDataContractPIVersion = "1"
)

type BlobDataContract struct {
	ResourceType                  string    `json:"resourceType" validate:"required"`
	ID                            string    `json:"id,omitempty"`
	Name                          string    `json:"name" validate:"required,min=1,max=20"`
	DataTypeID                    Reference `json:"dataTypeId" validate:"required,dive"`
	BucketID                      Reference `json:"bucketId" validate:"required,dive"`
	StorageClassID                Reference `json:"storageClassId" validate:"required,dive"`
	RootPathInBucket              string    `json:"rootPathInBucket" validate:"required,max=256"`
	LoggingEnabled                bool      `json:"loggingEnabled" validate:"required"`
	CrossRegionReplicationEnabled bool      `json:"crossRegionReplicationEnabled"`
	Meta                          *Meta     `json:"meta,omitempty"`
}

// GetBlobDataContractOptions struct describes search criteria for looking up BlobDataContract
type GetBlobDataContractOptions struct {
	ID   *string `url:"_id,omitempty"`
	Name *string `url:"name,omitempty"`
}

// Create creates a BlobDataContract
func (c *BlobDataContractsService) Create(ac BlobDataContract) (*BlobDataContract, *Response, error) {
	ac.ResourceType = "BlobDataContract"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}

	req, _ := c.NewRequest(http.MethodPost, "/BlobDataContract", ac, nil)
	req.Header.Set("api-version", blobDataContractPIVersion)

	var created BlobDataContract

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

// Delete deletes the given ServiceAction
func (c *BlobDataContractsService) Delete(ac BlobDataContract) (bool, *Response, error) {
	req, err := c.NewRequest(http.MethodDelete, "/BlobDataContract/"+ac.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobDataContractPIVersion)

	var deleteResponse interface{}

	resp, err := c.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// GetByID finds a client by its ID
func (c *BlobDataContractsService) GetByID(id string) (*BlobDataContract, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetById: missing id")
	}
	req, err := c.NewRequest(http.MethodGet, "/BlobDataContract/"+id, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobDataContractPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource BlobDataContract

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

// Find looks up services based on GetServiceActionOptions
func (c *BlobDataContractsService) Find(opt *GetBlobDataContractOptions, options ...OptionFunc) (*[]BlobDataContract, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, "/BlobDataContract", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobDataContractPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var resources []BlobDataContract
	for _, c := range bundleResponse.Entry {
		var resource BlobDataContract
		if err := json.Unmarshal(c.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, err
}

// Update updates a standard service
func (c *BlobDataContractsService) Update(ac BlobDataContract) (*BlobDataContract, *Response, error) {
	ac.ResourceType = "BlobDataContract"
	if err := c.validate.Struct(ac); err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest(http.MethodPut, "/BlobDataContract/"+ac.ID, ac, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobDataContractPIVersion)

	var updated BlobDataContract

	resp, err := c.Do(req, &updated)
	if err != nil {
		return nil, resp, err
	}
	return &updated, resp, nil
}
