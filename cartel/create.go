package cartel

type CreateResponse struct {
	Message []struct {
		EipAddress interface{} `json:"eip_address"`
		InstanceID string      `json:"instance_id"`
		IPAddress  string      `json:"ip_address"`
		Name       string      `json:"name"`
		Role       string      `json:"role"`
	} `json:"message,omitempty"`
	Result      string `json:"result,omitempty"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (cr CreateResponse) Success() bool {
	return cr.Result == "Success"
}

func (cr CreateResponse) InstanceID() string {
	if len(cr.Message) == 0 {
		return ""
	}
	return cr.Message[0].InstanceID
}

func (cr CreateResponse) IPAddress() string {
	if len(cr.Message) == 0 {
		return ""
	}
	return cr.Message[0].IPAddress
}

func (c *Client) Create(tagName string, opts ...RequestOptionFunc) (*CreateResponse, *Response, error) {
	var body RequestBody
	body.NameTag = []string{tagName}
	if body.Role == "" {
		body.Role = "container-host"
	}

	for _, f := range opts {
		if f != nil {
			if err := f(&body); err != nil {
				return nil, nil, err
			}
		}
	}
	req, err := c.NewRequest("POST", "v3/api/create", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody CreateResponse

	resp, err := c.Do(req, &responseBody)
	return &responseBody, resp, err
}
