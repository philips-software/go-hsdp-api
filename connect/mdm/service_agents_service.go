package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type ServiceAgentsService struct {
	*Client
}

type ServiceAgent struct {
	ResourceType            string    `json:"resourceType"`
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	Configuration           string    `json:"configuration"`
	DataSubscriberId        Reference `json:"dataSubscriberId"`
	AuthenticationMethodIds []string  `json:"authenticationMethodIds"`
	APIVersionSupported     string    `json:"apiVersionSupported"`
	Meta                    *Meta     `json:"meta"`
}

type GetServiceAgentOptions struct {
	LastUpdate       *string `url:"_lastUpdated,omitempty"`
	ReverseInclude   *string `url:"_revinclude,omitempty"`
	Include          *string `url:"_include,omitempty"`
	ID               *string `url:"_id,omitempty"`
	Name             *string `url:"name,omitempty"`
	DataSubscriberID *string `url:"dataSubscriberId,omitempty"`
}

func (r *ServiceAgentsService) Get(opt *GetServiceAgentOptions) (*[]ServiceAgent, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/ServiceAgent", opt)
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
	var resources []ServiceAgent
	for _, s := range bundleResponse.Entry {
		var resource ServiceAgent
		if err := json.Unmarshal(s.Resource, &resource); err == nil {
			resources = append(resources, resource)
		}
	}
	return &resources, resp, nil
}

func (r *ServiceAgentsService) GetByID(id string) (*ServiceAgent, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetByID: missing id")
	}
	resources, resp, err := r.Get(&GetServiceAgentOptions{
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
