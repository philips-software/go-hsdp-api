package cartel

func (c *Client) GetSecurityGroups() (*[]string, *Response, error) {
	var body RequestBody

	req, err := c.NewRequest("POST", "v3/api/get_security_groups", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody []string
	resp, err := c.Do(req, &responseBody)
	return &responseBody, resp, err
}
