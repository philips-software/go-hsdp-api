package cce

import (
	"fmt"
	"net/http"
)

type API struct {
	ID      string   `json:"id"`
	Version string   `json:"version"`
	BaseURL string   `json:"baseUrl"`
	Options []string `json:"options"`
}

type Role struct {
	Pattern string `json:"pattern"`
	Version string `json:"version"`
	Role    string `json:"role"`
	Apis    []API  `json:"apis"`
}

type Application struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

type Carehub struct {
	Carehub struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"carehub"`
	Self         Application   `json:"self"`
	Applications []Application `json:"applications"`
}

type DiscoveryResponse struct {
	Exp      int       `json:"exp"`
	Carehubs []Carehub `json:"carehubs"`
}

func (c *Client) Discovery() (*DiscoveryResponse, *Response, error) {
	req, err := c.NewRequest("GET", c.Endpoints.DiscoveryEndpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var discoveryResponse DiscoveryResponse
	resp, err := c.Do(req, &discoveryResponse)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return &discoveryResponse, resp, nil
}
