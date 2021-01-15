package cartel

func (c *Client) RemoveSecurityGroups(instances []string, groups []string) (*SecurityGroupsResponse, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.SecurityGroup = groups

	req, err := c.newRequest("POST", "v3/api/remove_security_groups", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody SecurityGroupsResponse
	resp, err := c.do(req, &responseBody)

	return &responseBody, resp, err
}
