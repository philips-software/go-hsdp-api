package cartel

import (
	"encoding/json"
	"net/http"
)

type StartResponse struct {
	Message     json.RawMessage `json:"message,omitempty"`
	Code        int             `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (sr StartResponse) Success() bool {
	return sr.Code == http.StatusOK
}

func (c *Client) Start(nameTag string) (*StartResponse, *Response, error) {
	var body RequestBody
	body.NameTag = []string{nameTag}

	req, err := c.newRequest("POST", "v3/api/start", &body, nil)
	if err != nil {
		return nil, nil, err
	}
	var responseBody StartResponse
	resp, err := c.do(req, &responseBody)
	if resp != nil {
		responseBody.Code = resp.StatusCode
	}
	return &responseBody, resp, err
}
