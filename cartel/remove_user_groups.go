package cartel

func (c *Client) RemoveUserGroups(instances []string, groups []string) (*UserGroupsResponse, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.LDAPGroups = groups

	req, err := c.newRequest("POST", "v3/api/remove_ldap_group", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody UserGroupsResponse
	resp, err := c.do(req, &responseBody)

	return &responseBody, resp, err
}
