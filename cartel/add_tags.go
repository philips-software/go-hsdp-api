package cartel

type AddTagResponse struct {
	Message     string `json:"message,omitempty"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (atr AddTagResponse) Success() bool {
	return atr.Code == 0
}

func (c *Client) AddTags(instances []string, tags map[string]string) (*AddTagResponse, *Response, error) {
	var body RequestBody
	body.NameTag = instances
	body.Tags = tags

	req, err := c.newRequest("POST", "v3/api/add_tags", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody AddTagResponse
	resp, err := c.do(req, &responseBody)

	return &responseBody, resp, err

}
