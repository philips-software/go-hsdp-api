package cartel

func (c *Client) AddTags(instances []string, tags map[string]string) (interface{}, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.Tags = tags

	req, err := c.NewRequest("POST", "v3/api/add_tags", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody interface{}
	resp, err := c.Do(req, &responseBody)

	return &responseBody, resp, err

}
