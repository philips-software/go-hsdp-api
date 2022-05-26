package blr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/internal"
)

type BlobsService struct {
	*Client
	validate *validator.Validate
}

var (
	blobAPIVersion = "1"
)

type Blob struct {
	ResourceType             string      `json:"resourceType" validate:"required"`
	ID                       string      `json:"id,omitempty"`
	DataType                 string      `json:"dataType" validate:"required"`
	Guid                     string      `json:"guid,omitempty"`
	Tags                     *[]Tag      `json:"tags,omitempty" validate:"omitempty,max=10"`
	AutoGenerateBlobPathName bool        `json:"autoGenerateBlobPathName"`
	BlobPath                 string      `json:"blobPath,omitempty" validate:"omitempty"`
	BlobName                 string      `json:"blobName,omitempty" validate:"omitempty"`
	VirtualPath              string      `json:"virtualPath,omitempty" validate:"omitempty"`
	VirtualName              string      `json:"virtualName,omitempty" validate:"omitempty"`
	Bucket                   string      `json:"bucket,omitempty"`
	Creation                 *string     `json:"creation,omitempty"`
	CreatedBy                string      `json:"createdBy,omitempty"`
	Attachment               *Attachment `json:"attachment,omitempty"`
	UploadOnBehalf           bool        `json:"uploadOnBehalf"`
	ManagingOrganization     string      `json:"managingOrganization,omitempty"`
	PropositionGuid          string      `json:"propositionGuid,omitempty"`
	MultipartEnabled         bool        `json:"multipartEnabled"`
	NoOfParts                *int        `json:"noOfParts,omitempty"`
	State                    *string     `json:"state,omitempty"`
	Meta                     *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	VersionID   string    `json:"versionId,omitempty"`
}

type Tag struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type Attachment struct {
	ContentType string `json:"contentType,omitempty"`
	Language    string `json:"language,omitempty"`
	Hash        string `json:"hash,omitempty"`
	Title       string `json:"title,omitempty"`
	Data        string `json:"data" validate:"required"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	Created     string `json:"created,omitempty"`
}

type GetBlobOptions struct {
	HSDPID          *string `url:"hsdpId,omitempty"`
	BlobName        *string `url:"blobName,omitempty"`
	DataType        *string `url:"dataType,omitempty"`
	GUID            *string `url:"guid,omitempty"`
	DefaultRegionID *string `url:"defaultRegionId,omitempty"`
	LastUpdated     *string `url:"_lastUpdated,omitempty"`
	StartDate       *string `url:"startDate,omitempty"`
	EndDate         *string `url:"endDate,omitempty"`
	Include         *string `url:"_include,omitempty"`
	Page            *string `url:"page,omitempty"`
	Count           *int    `url:"_count,omitempty"`
	SingleDownload  *string `url:"singleDownload,omitempty"`
}

type BlobPolicy struct {
	ResourceType string            `json:"resourceType" validate:"required"`
	Statement    []PolicyStatement `json:"statement" validate:"required"`
}

type PolicyStatement struct {
	ResourceType string                `json:"resourceType" validate:"required"`
	SID          *string               `json:"sid,omitempty"`
	Principal    PrincipalResourceList `json:"principal" validate:"required"`
	Effect       string                `json:"effect" validate:"required"`
	Action       []string              `json:"action" validate:"required"`
}

type PrincipalResourceList struct {
	HSDP []string `json:"hsdp" validate:"required,min=1,max=64"`
}

type AccessURL struct {
	ResourceType string     `json:"resourceType"`
	Actions      []string   `json:"actions"`
	URL          string     `json:"url"`
	URLExpiry    string     `json:"urlExpiry"`
	BlobPartURLs []BlobPart `json:"blobPartUrls,omitempty"`
}

type BlobPart struct {
	PartNumber          int    `json:"partNumber"`
	DataAccessURL       string `json:"dataAccessUrl"`
	DataAccessUrlExpiry string `json:"dataAccessUrlExpiry"`
}

type BlobPartUpload struct {
	ResourceType string       `json:"resourceType"`
	BlobParts    []PartUpload `json:"blobParts"`
}

type PartUpload struct {
	PartNumber   int    `json:"partNumber"`
	Size         int    `json:"size"`
	LastModified string `json:"lastModified"`
	ETag         string `json:"eTag"`
}

func (b *BlobsService) Create(blob Blob) (*Blob, *Response, error) {
	blob.ResourceType = "Blob"
	blob.AutoGenerateBlobPathName = true
	if err := b.validate.Struct(blob); err != nil {
		return nil, nil, err
	}

	req, _ := b.NewRequest(http.MethodPost, "/Blob", blob, nil)
	req.Header.Set("api-version", blobAPIVersion)

	var created Blob

	resp, err := b.Do(req, &created)

	if err != nil {
		return nil, resp, err
	}
	if created.ID == "" {
		return nil, resp, fmt.Errorf("the 'ID' field is missing")
	}
	return &created, resp, nil
}

func (b *BlobsService) GetByID(id string) (*Blob, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Blob/"+id, nil)
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

func (b *BlobsService) Find(opt *GetBlobOptions, options ...OptionFunc) (*[]Blob, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Blob", opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
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

func (b *BlobsService) Delete(blob Blob) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/Blob/"+blob.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)

	var deleteResponse interface{}

	resp, err := b.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

func (b *BlobsService) SetPolicy(blob Blob, policy BlobPolicy) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodPost, "/Blob/"+blob.ID+"/$setPolicy", policy, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)

	var setPolicyResponse interface{}

	resp, err := b.Do(req, &setPolicyResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

func (b *BlobsService) GetPolicy(blob Blob) (*BlobPolicy, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Blob/"+blob.ID+"/$getPolicy", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource BlobPolicy

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetPolicy: %w", err)
	}
	return &resource, resp, nil
}

func (b *BlobsService) DeletePolicy(blob Blob) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodDelete, "/Blob/"+blob.ID+"/$deletePolicy", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)

	var deleteResponse interface{}

	resp, err := b.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

func (b *BlobsService) GetAccessURL(blob Blob) (*AccessURL, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Blob/"+blob.ID+"/$getAccessUrl", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource AccessURL

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetAccessURL: %w", err)
	}
	return &resource, resp, nil
}

func (b *BlobsService) CompleteUpload(blob Blob, parts BlobPartUpload) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodPost, "/Blob/"+blob.ID+"/$completeUpload", parts, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)

	var completeUploadResponse interface{}

	resp, err := b.Do(req, &completeUploadResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

func (b *BlobsService) AbortUpload(blob Blob) (bool, *Response, error) {
	req, err := b.NewRequest(http.MethodPost, "/Blob/"+blob.ID+"/$abortUpload", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)

	var abortUploadResponse interface{}

	resp, err := b.Do(req, &abortUploadResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

func (b *BlobsService) ListParts(blob Blob) (*BlobPartUpload, *Response, error) {
	req, err := b.NewRequest(http.MethodGet, "/Blob/"+blob.ID+"/$listPart", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", blobAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var resource BlobPartUpload

	resp, err := b.Do(req, &resource)
	if err != nil {
		return nil, resp, err
	}
	err = internal.CheckResponse(resp.Response)
	if err != nil {
		return nil, resp, fmt.Errorf("GetByID: %w", err)
	}
	return &resource, resp, nil
}
