package cartel

func (c *Client) RemoveUserGroups(instances []string, groups []string) (*UserGroupsResponse, *Response, error) {
	var body RequestBody
	var responseBody UserGroupsResponse
	var resp *Response
	var err error
	body.NameTag = instances

	for _, group := range groups {
		body.LDAPGroups = []string{group}
		req, err := c.newRequest("POST", "v3/api/remove_ldap_group", &body, nil)
		if err != nil {
			return nil, nil, err
		}
		resp, err = c.do(req, &responseBody)
	}
	return &responseBody, resp, err
}
