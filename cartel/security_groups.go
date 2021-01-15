package cartel

func (c *Client) GetSecurityGroups() (*[]string, *Response, error) {
	var body RequestBody

	req, err := c.newRequest("POST", "v3/api/get_security_groups", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody []string
	resp, err := c.do(req, &responseBody)
	return &responseBody, resp, err
}
