package cartel

type Role struct {
	Description string `json:"description"`
	Role        string `json:"role"`
}

func (c *Client) GetRoles() (*[]Role, *Response, error) {
	var body RequestBody
	body.Token = c.config.Token

	req, err := c.newRequest("POST", "v3/api/get_all_roles", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var roleResponse []Role

	resp, err := c.do(req, &roleResponse)
	return &roleResponse, resp, err
}
