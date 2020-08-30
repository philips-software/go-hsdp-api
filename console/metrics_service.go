package console

import "time"

type MetricsService struct {
	client *Client
}

type Instance struct {
	CreatedAt    time.Time `json:"createdAt"`
	GUID         string    `json:"guid"`
	Name         string    `json:"name"`
	Organization string    `json:"organization"`
	Space        string    `json:"space"`
}

type MetricsResponse struct {
	Data struct {
		Instances []Instance `json:"instances"`
	} `json:"data"`
	Status string `json:"status"`
}

// GetInstances looks up available instances
func (c *MetricsService) GetInstances(options ...OptionFunc) (*[]Instance, *Response, error) {
	req, err := c.client.NewRequest(CONSOLE, "GET", "v3/metrics/instances", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response MetricsResponse

	resp, err := c.client.Do(req, &response)
	if err != nil {
		return nil, resp, err
	}
	return &response.Data.Instances, resp, err
}

// GetInstances looks up available instances
func (c *MetricsService) GetInstanceByID(id string, options ...OptionFunc) (*Instance, *Response, error) {
	req, err := c.client.NewRequest(CONSOLE, "GET", "v3/metrics/intances/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response struct {
		Data   Instance `json:"data"`
		Status string   `json:"status"`
	}

	resp, err := c.client.Do(req, &response)
	if err != nil {
		return nil, resp, err
	}
	return &response.Data, resp, err
}
