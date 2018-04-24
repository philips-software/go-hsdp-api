package api

import (
	"io/ioutil"
	"net/url"
	"strings"
)

type IntrospectResponse struct {
	Active        bool   `json:"active"`
	Scope         string `json:"scope"`
	ISS           string `json:"iss"`
	Username      string `json:"username"`
	Expires       int64  `json:"exp"`
	Organizations struct {
		ManagingOrganization string `json:"managingOrganization"`
		OrganizationList     []struct {
			OrganizationID string   `json:"organizationId"`
			Permissions    []string `json:"permissions"`
		} `json:"organizationList"`
	} `json:"organizations"`
	ClientID     string `json:"client_id"`
	IdentityType string `json:"identity_type"`
}

func (c *Client) Introspect() (*IntrospectResponse, *Response, error) {
	var val IntrospectResponse

	req, err := c.NewIAMRequest("POST", "authorize/oauth2/introspect", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	form := url.Values{}
	form.Add("token", c.token)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req, &val)

	return &val, resp, err
}
