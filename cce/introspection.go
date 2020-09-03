package cce

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	Exp       int    `json:"exp"`
	Sub       string `json:"sub"`
	Iss       string `json:"iss"`
	TokenType string `json:"token_type"`
	Aud       string `json:"aud"`
	ClientID  string `json:"client_id"`
	CarehubID string `json:"carehub_id"`
}

func (c *Client) Introspection() (*IntrospectionResponse, *Response, error) {
	req, err := c.NewRequest("POST", c.Endpoints.IntrospectionEndpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	form := url.Values{}
	form.Add("token", c.iamClient.Token())
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Api-Version", APIVersion)

	var val IntrospectionResponse
	resp, err := c.Do(req, &val)

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return &val, resp, nil
}
