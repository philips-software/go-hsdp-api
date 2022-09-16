package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

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
	Error  Error  `json:"error,omitempty"`
}

type Group struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type RuleResponse struct {
	Data struct {
		Groups []Group `json:"groups"`
	} `json:"data"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Threshold struct {
	Default int      `json:"default,omitempty"`
	Enabled bool     `json:"enabled"`
	Max     float64  `json:"max"`
	Min     float64  `json:"min"`
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
	Unit    []string `json:"unit,omitempty"`
}

type Application struct {
	Enabled      bool        `json:"enabled"`
	MaxInstances int         `json:"maxInstances"`
	MinInstances int         `json:"minInstances"`
	Name         string      `json:"name"`
	Thresholds   []Threshold `json:"thresholds,omitempty"`
}

type AutoscalersResponse struct {
	Data struct {
		Applications []Application `json:"applications"`
	} `json:"data"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Rule struct {
	Annotations struct {
		Description string `json:"description"`
		Resolved    string `json:"resolved"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Metric      string `json:"metric"`
	Rule        struct {
		ExtraFor []struct {
			Name         string   `json:"name"`
			Options      []string `json:"options"`
			Type         string   `json:"type"`
			VariableName string   `json:"variableName"`
		} `json:"extraFor,omitempty"`
		Extras []struct {
			Name         string   `json:"name"`
			Options      []string `json:"options"`
			Type         string   `json:"type"`
			VariableName string   `json:"variableName"`
		} `json:"extras"`
		Operators []string  `json:"operators"`
		Subject   string    `json:"subject"`
		Threshold Threshold `json:"threshold"`
	} `json:"rule"`
	Template string `json:"template"`
}

// GetGroupedRules looks up available rules
func (c *MetricsService) GetGroupedRules(options ...OptionFunc) (*[]Group, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/rules", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var jsonResponse RuleResponse
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&jsonResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", jsonResponse.Status, jsonResponse.Error.Code, jsonResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &jsonResponse.Data.Groups, resp, err
}

// GetRuleByID retrieves a rule by ID
func (c *MetricsService) GetRuleByID(id string, options ...OptionFunc) (*Rule, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/rules/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var jsonResponse struct {
		Data   Rule   `json:"data"`
		Status string `json:"status"`
		Error  Error  `json:"error,omitempty"`
	}
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&jsonResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", jsonResponse.Status, jsonResponse.Error.Code, jsonResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &jsonResponse.Data, resp, err
}

// GetInstances looks up available instances
func (c *MetricsService) GetInstances(options ...OptionFunc) (*[]Instance, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/instances", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var jsonResponse MetricsResponse
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&jsonResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", jsonResponse.Status, jsonResponse.Error.Code, jsonResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &jsonResponse.Data.Instances, resp, err
}

// GetInstanceByID looks up an instance by ID
func (c *MetricsService) GetInstanceByID(id string, options ...OptionFunc) (*Instance, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/instances/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var jsonResponse struct {
		Data   Instance `json:"data"`
		Status string   `json:"status"`
		Error  Error    `json:"error,omitempty"`
	}
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&jsonResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", jsonResponse.Status, jsonResponse.Error.Code, jsonResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &jsonResponse.Data, resp, err
}

// GetApplicationAutoscalers looks up all available autoscalers
func (c *MetricsService) GetApplicationAutoscalers(id string, options ...OptionFunc) (*[]Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/"+id+"/autoscalers", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var jsonResponse AutoscalersResponse
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&jsonResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", jsonResponse.Status, jsonResponse.Error.Code, jsonResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &jsonResponse.Data.Applications, resp, err
}

// GetApplicationAutoscaler looks up a specific application autoscaler settings
func (c *MetricsService) GetApplicationAutoscaler(id, app string, options ...OptionFunc) (*Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/"+id+"/autoscalers/"+app, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var getResponse struct {
		Data struct {
			Application Application `json:"application"`
		} `json:"data"`
		Status string `json:"status"`
		Error  Error  `json:"error,omitempty"`
	}
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&getResponse)
	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", getResponse.Status, getResponse.Error.Code, getResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &getResponse.Data.Application, resp, err
}

// UpdateApplicationAutoscaler updates a specific application autoscaler settings
func (c *MetricsService) UpdateApplicationAutoscaler(id string, settings Application, options ...OptionFunc) (*Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "PUT", "v3/metrics/"+id+"/autoscalers", &settings, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var updateResponse struct {
		Data struct {
			Application Application `json:"application"`
		} `json:"data,omitempty"`
		Status string `json:"status,omitempty"`
		Error  Error  `json:"error,omitempty"`
	}
	var response bytes.Buffer

	resp, err := c.client.do(req, &response)
	jsonErr := json.NewDecoder(&response).Decode(&updateResponse)

	if err != nil {
		if jsonErr == nil {
			return nil, resp, fmt.Errorf("status: %s, code: %s, message: %s, error: %w", updateResponse.Status, updateResponse.Error.Code, updateResponse.Error.Message, err)
		}
		return nil, resp, err
	}
	if jsonErr != nil {
		return nil, resp, fmt.Errorf("decoding jsonResponse: %w", err)
	}
	return &updateResponse.Data.Application, resp, err
}
