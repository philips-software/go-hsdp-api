package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type SubscriberTypesService struct {
	*Client
}

type SubscriberType struct {
	ResourceType          string `json:"resourceType"`
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	ConfigurationTemplate string `json:"configurationTemplate"`
	SubscriptionTemplate  string `json:"subscriptionTemplate"`
	Meta                  *Meta  `json:"meta,omitempty"`
}

type GetSubscriberTypeOptions struct {
	LastUpdate     *string `url:"_lastUpdated,omitempty"`
	ReverseInclude *string `url:"_revinclude,omitempty"`
	Include        *string `url:"_include,omitempty"`
	ID             *string `url:"_id,omitempty"`
	Name           *string `url:"name,omitempty"`
}

func (r *SubscriberTypesService) Get(opt *GetSubscriberTypeOptions) (*[]SubscriberType, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/SubscriberType", opt)
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
	var resources []SubscriberType
	for _, s := range bundleResponse.Entry {
		var resource SubscriberType
		if err := json.Unmarshal(s.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, nil
}

func (r *SubscriberTypesService) GetByID(id string) (*SubscriberType, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetByID: missing id")
	}
	resources, resp, err := r.Get(&GetSubscriberTypeOptions{
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
