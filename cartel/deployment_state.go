package cartel

func (c *Client) GetDeploymentState(nameTag string) (string, *Response, error) {
	var body RequestBody
	body.NameTag = []string{nameTag}

	req, err := c.NewRequest("POST", "v3/api/deployment_status", &body, nil)
	if err != nil {
		return "fatal_error", nil, err
	}
	var responseBody map[string]interface{}
	resp, err := c.Do(req, &responseBody)
	if err != nil {
		return "unknown_instance", resp, err
	}
	state, ok := responseBody[nameTag].(map[string]interface{})
	if !ok {
		return "unknown_instance", resp, err
	}
	deployState, ok := state["deploy_state"].(string)
	if !ok {
		return "indeterminate", resp, err
	}
	return deployState, resp, err
}
