package cartel

type SecurityGroupDetails []SecurityRule

type SecurityRule struct {
	PortRange string   `json:"port_range"`
	Protocol  string   `json:"protocol"`
	Source    []string `json:"source"`
}

func (c *Client) GetSecurityGroupDetails(group string) (*SecurityGroupDetails, *Response, error) {
	var body RequestBody
	body.SecurityGroup = []string{group}

	req, err := c.NewRequest("POST", "v3/api/security_group_details", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody map[string]SecurityGroupDetails
	resp, err := c.Do(req, &responseBody)
	details := responseBody[group]
	return &details, resp, err
}
