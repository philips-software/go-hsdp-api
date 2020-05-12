package cartel

type DestroyResponse struct {
	AWS    string            `json:"AWS"`
	Cartel map[string]string `json:"Cartel"`
}

func (dr DestroyResponse) Success() bool {
	// This works since we only handle one instance per Create/Destroy call
	for _, v := range dr.Cartel {
		if v == "Instance removed." {
			return true
		}
	}
	return false
}

func (c *Client) DestroyInstance(tagName string) (*DestroyResponse, *Response, error) {
	var body RequestBody
	body.NameTag = []string{tagName}

	req, err := c.NewRequest("POST", "v3/api/destroy", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody DestroyResponse

	resp, err := c.Do(req, &responseBody)
	if err != nil {
		return nil, resp, err
	}
	return &responseBody, resp, err
}
