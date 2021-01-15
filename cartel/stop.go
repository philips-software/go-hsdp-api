package cartel

import "encoding/json"

type StopResponse struct {
	Message     json.RawMessage `json:"message,omitempty"`
	Code        int             `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (sr StopResponse) Success() bool {
	return sr.Code == 0
}

func (c *Client) Stop(nameTag string) (*StopResponse, *Response, error) {
	var body RequestBody
	body.NameTag = []string{nameTag}

	req, err := c.newRequest("POST", "v3/api/suspend", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody StopResponse
	resp, err := c.do(req, &responseBody)
	return &responseBody, resp, err
}
