package stl

import (
	"context"
	"fmt"
	"github.com/hasura/go-graphql-client"
)

type ConfigService struct {
	client *Client
}

type AppFirewallException struct {
	DeviceID int64 `json:"deviceId,omitempty"`
	TCP      []int `json:"tcp"`
	UDP      []int `json:"udp"`
}

type AppLogging struct {
	DeviceID         int64  `json:"deviceId,omitempty"`
	RawConfig        string `json:"rawConfig"`
	HSDPLogging      bool   `json:"hsdpLogging"`
	HSDPIngestorHost string `json:"hsdpIngestorHost"`
	HSDPSharedKey    string `json:"hsdpSharedKey"`
	HSDPSecretKey    string `json:"hsdpSecretKey"`
	HSDPProductKey   string `json:"hsdpProductKey"`
	HSDPCustomField  *bool  `json:"hsdpCustomField,omitempty"`
}

type UpdateAppFirewallExceptionInput struct {
	AppFirewallException
	SerialNumber string `json:"serialNumber"`
}

type UpdateAppLoggingInput struct {
	AppLogging
	SerialNumber string `json:"serialNumber"`
}

func (u UpdateAppLoggingInput) Validate() (bool, error) {
	// Empty config is acceptable
	if u.RawConfig == "" && u.HSDPSharedKey == "" && u.HSDPIngestorHost == "" && u.HSDPSecretKey == "" && u.HSDPProductKey == "" {
		return true, nil
	}
	// Valid HSDP config
	if u.HSDPSharedKey != "" && u.HSDPIngestorHost != "" && u.HSDPSecretKey != "" && u.HSDPProductKey != "" && u.RawConfig == "" {
		return true, nil
	}
	// Valid RAW config
	if u.RawConfig != "" && (u.HSDPSharedKey == "" && u.HSDPIngestorHost == "" && u.HSDPSecretKey == "" && u.HSDPProductKey == "") {
		return true, nil
	}
	return false, fmt.Errorf("invalid or conflicting HSDP / Raw config")
}

func (c *ConfigService) GetFirewallExceptionsBySerial(ctx context.Context, serial string) (*AppFirewallException, error) {
	var query struct {
		AppFirewallException AppFirewallException `graphql:"appFirewallException(serialNumber: $serialNumber)"`
	}
	err := c.client.gql.Query(ctx, &query, map[string]interface{}{
		"serialNumber": graphql.String(serial),
	})
	if err != nil {
		return nil, err
	}
	return &query.AppFirewallException, nil
}

func (c *ConfigService) UpdateAppFirewallExceptions(ctx context.Context, input UpdateAppFirewallExceptionInput) (*AppFirewallException, error) {
	var mutation struct {
		UpdateAppFirewallException struct {
			StatusCode           int
			Success              bool
			Message              string
			AppFirewallException AppFirewallException
		} `graphql:"updateAppFirewallException(input: $input)"`
	}
	err := c.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.UpdateAppFirewallException.AppFirewallException, nil
}

func (c *ConfigService) GetAppLoggingBySerial(ctx context.Context, serial string) (*AppLogging, error) {
	var query struct {
		AppLogging AppLogging `graphql:"appLogging(serialNumber: $serialNumber)"`
	}
	err := c.client.gql.Query(ctx, &query, map[string]interface{}{
		"serialNumber": graphql.String(serial),
	})
	if err != nil {
		return nil, err
	}
	return &query.AppLogging, nil
}

func (c *ConfigService) UpdateAppLogging(ctx context.Context, input UpdateAppLoggingInput) (*AppLogging, error) {
	var mutation struct {
		UpdateAppLogging struct {
			StatusCode int
			Success    bool
			Message    string
			AppLogging AppLogging
		} `graphql:"updateAppLogging(input: $input)"`
	}
	err := c.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.UpdateAppLogging.AppLogging, nil
}
