package mdm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type OAuthClientScopesService struct {
	*Client
}

type OAuthClientScope struct {
	ResourceType     string `json:"resourceType" validate:"required"`
	ID               string `json:"id"`
	Organization     string `json:"organization"`
	Proposition      string `json:"proposition"`
	Service          string `json:"service"`
	Resource         string `json:"resource"`
	Action           string `json:"action"`
	BootstrapEnabled bool   `json:"bootstrapEnabled"`
	Meta             *Meta  `json:"meta"`
}

type GetOAuthClientScopeOptions struct {
	ID           *string `url:"_id,omitempty"`
	LastUpdate   *string `url:"_lastUpdated,omitempty"`
	Organization *string `url:"organization,omitempty"`
	Proposition  *string `url:"proposition,omitempty"`
	Action       *string `url:"service,omitempty"`
	Service      *string `url:"action,omitempty"`
}

func (r *OAuthClientScopesService) GetOAuthClientScopes(opt *GetOAuthClientScopeOptions) (*[]OAuthClientScope, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/OAuthClientScope", opt)
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
	var scopes []OAuthClientScope
	for _, s := range bundleResponse.Entry {
		var scope OAuthClientScope
		err := json.Unmarshal(s.Resource, &scope)
		if err == nil {
			scopes = append(scopes, scope)
		}
	}
	return &scopes, resp, nil
}

func (r *OAuthClientScopesService) GetOAuthClientScopeByID(id string) (*OAuthClientScope, *Response, error) {
	if len(id) == 0 {
		return nil, nil, fmt.Errorf("GetOAuthClientScopeByID: missing id")
	}
	classes, resp, err := r.GetOAuthClientScopes(&GetOAuthClientScopeOptions{
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
