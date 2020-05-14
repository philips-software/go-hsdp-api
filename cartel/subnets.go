package cartel

type Subnet struct {
	ID      string `json:"id"`
	Network string `json:"network"`
}

type SubnetDetails map[string]Subnet

func (c *Client) GetAllSubnets() (*SubnetDetails, *Response, error) {
	var body RequestBody

	req, err := c.NewRequest("POST", "v3/api/get_all_subnets", &body, nil)
	if err != nil {
		return nil, nil, err
	}

	var responseBody SubnetDetails

	resp, err := c.Do(req, &responseBody)

	return &responseBody, resp, err
}
