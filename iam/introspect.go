package iam

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	introspectAPIVersion = "3"
)

// IntrospectResponse contains details of the introspect on a profile
type IntrospectResponse struct {
	Active        bool   `json:"active"`
	Scope         string `json:"scope"`
	Username      string `json:"username"`
	Expires       int64  `json:"exp"`
	Sub           string `json:"sub"`
	ISS           string `json:"iss"`
	Organizations struct {
		ManagingOrganization string `json:"managingOrganization"`
		OrganizationList     []struct {
			OrganizationID   string   `json:"organizationId"`
			Permissions      []string `json:"permissions"`
			OrganizationName string   `json:"organizationName"`
			Groups           []string `json:"groups"`
			Roles            []string `json:"roles"`
		} `json:"organizationList"`
	} `json:"organizations"`
	ClientID     string `json:"client_id"`
	TokenType    string `json:"token_type"`
	IdentityType string `json:"identity_type"`
}

func WithOrgContext(organizationId string) OptionFunc {
	return func(req *http.Request) error {
		params := struct {
			OrgContext *string `url:"org_ctx"`
		}{
			OrgContext: &organizationId,
		}
		val, err := query.Values(params)
		if err != nil {
			return err
		}
		req.URL.RawQuery = val.Encode()
		return nil
	}
}

// Introspect introspects the current logged in user
func (c *Client) Introspect(opts ...OptionFunc) (*IntrospectResponse, *Response, error) {
	var val IntrospectResponse

	req, err := c.newRequest(IAM, "POST", "authorize/oauth2/introspect", nil, opts)
	if err != nil {
		return nil, nil, err
	}
	form := url.Values{}
	form.Add("token", c.token)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Api-Version", introspectAPIVersion)

	resp, err := c.do(req, &val)

	return &val, resp, err
}
