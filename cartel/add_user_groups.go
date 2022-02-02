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
	var responseBody UserGroupsResponse
	var resp *Response
	var err error
	body.NameTag = instances

	for _, group := range groups {
		body.LDAPGroups = []string{group} // Can only add/remove single group
		req, err := c.newRequest("POST", "v3/api/add_ldap_group", &body, nil)
		if err != nil {
			return nil, nil, err
		}
		resp, err = c.do(req, &responseBody)
	}
	return &responseBody, resp, err
}
