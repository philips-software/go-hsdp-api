package cartel

func (c *Client) GetAllInstances() (*[]InstanceDetails, *Response, error) {
	var body RequestBody

	req, err := c.newRequest("POST", "v3/api/get_all_instances", &body, nil)
	if err != nil {
		return nil, nil, err
	}

	var responseBody []InstanceDetails

	resp, err := c.do(req, &responseBody)

	return &responseBody, resp, err
}
