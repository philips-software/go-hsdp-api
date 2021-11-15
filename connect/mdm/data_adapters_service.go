package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type DataAdaptersService struct {
	*Client
}

type DataAdapter struct {
	ResourceType   string    `json:"resourceType"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceAgentId Reference `json:"serviceAgentId"`
	Meta           *Meta     `json:"meta"`
}

type GetDataAdapterOptions struct {
	LastUpdate          *string `url:"_lastUpdated,omitempty"`
	ReverseInclude      *string `url:"_revinclude,omitempty"`
	Include             *string `url:"_include,omitempty"`
	ID                  *string `url:"_id,omitempty"`
	Name                *string `url:"name,omitempty"`
	SubscriberTypeID    *string `url:"subscriberTypeId,omitempty"`
	SubscriberTypeFGUID *string `url:"subscriberGuid,omitempty"`
}

func (r *DataAdaptersService) Get(opt *GetDataAdapterOptions) (*[]DataAdapter, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/DataAdapter", opt)
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
	var resources []DataAdapter
	for _, s := range bundleResponse.Entry {
		var resource DataAdapter
		if err := json.Unmarshal(s.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, nil
}

func (r *DataAdaptersService) GetByID(id string) (*DataAdapter, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetByID: missing id")
	}
	resources, resp, err := r.Get(&GetDataAdapterOptions{
		ID: &id,
	})
	if err != nil {
		return nil, resp, err
	}
	if len(*resources) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &(*resources)[0], resp, nil
}
