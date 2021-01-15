package cartel

type ProtectionResponse struct {
	Message     string `json:"message,omitempty"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (pr ProtectionResponse) Success() bool {
	return pr.Code == 0
}

func (c *Client) SetProtection(nameTag string, protection bool) (*ProtectionResponse, *Response, error) {
	var body RequestBody
	body.NameTag = []string{nameTag}
	body.Protect = protection

	req, err := c.newRequest("POST", "v3/api/protect", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody ProtectionResponse
	resp, err := c.do(req, &responseBody)
	return &responseBody, resp, err
}
