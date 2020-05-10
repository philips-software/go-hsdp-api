package cartel

func (c *Client) GetRoles() (*Response, error) {
	var body CartelRequestBody
	body.Token = c.config.Token

	req, err := c.NewRequest("POST", "v3/api/get_all_roles", body, nil)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	return c.Do(req, &responseBody)
}
