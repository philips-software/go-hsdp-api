package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type DataSubscribersService struct {
	*Client
}

type DataSubscriber struct {
	ResourceType     string     `json:"resourceType"`
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	SubscriberGuid   Identifier `json:"subscriberGuid"`
	Configuration    string     `json:"configuration"`
	ApiVersion       string     `json:"apiVersion"`
	SubscriberTypeID string     `json:"subscriberTypeId"`
	Meta             *Meta      `json:"meta"`
}

type GetDataSubscriberOptions struct {
	LastUpdate          *string `url:"_lastUpdated,omitempty"`
	ReverseInclude      *string `url:"_revinclude,omitempty"`
	Include             *string `url:"_include,omitempty"`
	ID                  *string `url:"_id,omitempty"`
	Name                *string `url:"name,omitempty"`
	SubscriberTypeID    *string `url:"subscriberTypeId,omitempty"`
	SubscriberTypeFGUID *string `url:"subscriberGuid,omitempty"`
}

func (r *DataSubscribersService) Get(opt *GetDataSubscriberOptions) (*[]DataSubscriber, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/DataSubscriber", opt)
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
	var resources []DataSubscriber
	for _, s := range bundleResponse.Entry {
		var resource DataSubscriber
		if err := json.Unmarshal(s.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, nil
}

func (r *DataSubscribersService) GetByID(id string) (*DataSubscriber, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetByID: missing id")
	}
	resources, resp, err := r.Get(&GetDataSubscriberOptions{
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
