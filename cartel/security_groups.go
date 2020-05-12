package cartel

type SecurityGroupResponse struct {
}

func (c *Client) GetSecurityGroups() (*Response, error) {
	var body CartelRequestBody

	req, err := c.NewRequest("POST", "v3/api/get_security_groups", &body, nil)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	return c.Do(req, &responseBody)
}
