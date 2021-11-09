package mdm

import (
	"encoding/json"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type StorageClassService struct {
	*Client
}

type StorageClass struct {
	ResourceType string `json:"resourceType"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Meta         *Meta  `json:"meta"`
}

type GetStorageClassOptions struct {
	LastUpdate     *string `url:"_lastUpdated"`
	ReverseInclude *string `url:"_revinclude"`
	Include        *string `url:"_include"`
	ID             *string `url:"_id"`
	Name           *string `url:"name"`
}

func (r *StorageClassService) GetStorageClasses(opt *GetStorageClassOptions) (*[]StorageClass, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/Region", opt)
	if err != nil {
		return nil, nil, err
	}
	var bundleResponse internal.Bundle

	resp, err := r.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if err := internal.CheckResponse(resp.Response); err != nil {
		return nil, resp, err
	}
	var classes []StorageClass
	for _, s := range bundleResponse.Entry {
		var class StorageClass
		if err := json.Unmarshal(s.Resource, &class); err == nil {
			classes = append(classes, class)
		}
	}
	return &classes, resp, nil
}

func (r *StorageClassService) GetRegionByID(id string) (*StorageClass, *Response, error) {
	classes, resp, err := r.GetStorageClasses(&GetStorageClassOptions{
		ID: &id,
	})
	if err != nil {
		return nil, resp, err
	}
	if len(*classes) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &(*classes)[0], resp, nil
}
