package cartel

import (
	"encoding/json"
	"net/http"
)

type SecurityGroupsResponse struct {
	Message     json.RawMessage `json:"message,omitempty"`
	Result      string          `json:"result,omitempty"`
	Code        int             `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (sgr SecurityGroupsResponse) Success() bool {
	return sgr.Code == 0 || sgr.Code == http.StatusOK
}

func (c *Client) AddSecurityGroups(instances []string, groups []string) (*SecurityGroupsResponse, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.SecurityGroup = groups

	req, err := c.NewRequest("POST", "v3/api/add_security_groups", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody SecurityGroupsResponse
	resp, err := c.Do(req, &responseBody)

	return &responseBody, resp, err
}
