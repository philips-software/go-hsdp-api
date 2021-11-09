package mdm

import (
	"encoding/json"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type RegionsService struct {
	*Client
}

type Region struct {
	ResourceType string `json:"resourceType"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	HsdpEnabled  bool   `json:"hsdpEnabled"`
	Meta         *Meta  `json:"meta,omitempty"`
}

type GetRegionOptions struct {
	LastUpdate     *string `url:"_lastUpdated"`
	ReverseInclude *string `url:"_revinclude"`
	Include        *string `url:"_include"`
	ID             *string `url:"_id"`
	Name           *string `url:"name"`
	Category       *string `url:"category"`
	HSDPEnabled    *bool   `url:"hsdpEnabled"`
}

func (r *RegionsService) GetRegions(opt *GetRegionOptions) (*[]Region, *Response, error) {
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
	var regions []Region
	for _, s := range bundleResponse.Entry {
		var region Region
		if err := json.Unmarshal(s.Resource, &region); err == nil {
			regions = append(regions, region)
		}
	}
	return &regions, resp, nil
}

func (r *RegionsService) GetRegionByID(id string) (*Region, *Response, error) {
	regions, resp, err := r.GetRegions(&GetRegionOptions{
		ID: &id,
	})
	if err != nil {
		return nil, resp, err
	}
	if len(*regions) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &(*regions)[0], resp, nil
}
