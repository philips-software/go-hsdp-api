package cartel

func (c *Client) AddUserGroups(instances []string, groups []string) (interface{}, error) {
	var body RequestBody
	body.NameTag = instances
	body.LDAPGroups = groups

	req, err := c.NewRequest("POST", "v3/api/add_ldap_groups", &body, nil)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	return c.Do(req, &responseBody)
}
