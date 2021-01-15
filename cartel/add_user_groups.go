package cartel

import "encoding/json"

type UserGroupsResponse struct {
	Message     json.RawMessage `json:"message,omitempty"`
	Result      string          `json:"result,omitempty"`
	Code        int             `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (ugr UserGroupsResponse) Success() bool {
	return ugr.Code == 0
}

func (c *Client) AddUserGroups(instances []string, groups []string) (*UserGroupsResponse, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.LDAPGroups = groups

	req, err := c.newRequest("POST", "v3/api/add_ldap_group", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody UserGroupsResponse
	resp, err := c.do(req, &responseBody)

	return &responseBody, resp, err
}
