package iam

import (
	"io/ioutil"
	"net/url"
	"strings"
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

// Introspect introspects the current logged in user
func (c *Client) Introspect() (*IntrospectResponse, *Response, error) {
	var val IntrospectResponse

	req, err := c.newRequest(IAM, "POST", "authorize/oauth2/introspect", nil, nil)
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
