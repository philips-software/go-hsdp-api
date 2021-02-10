package stl

import (
	"context"
	"fmt"
	"github.com/hasura/go-graphql-client"
)

// Device represents a STL device
type Device struct {
	ID               int64
	Name             string
	State            string
	Region           string
	SerialNumber     string
	PrimaryInterface struct {
		Name    string
		Address string
	}
}

type DevicesService struct {
	client *Client
}

type SyncDeviceConfigsInput struct {
	SerialNumber string `json:"serialNumber"`
}

func (d *DevicesService) SyncDeviceConfig(ctx context.Context, serial string) error {
	var mutation struct {
		SyncDeviceConfigs struct {
			StatusCode int
			Success    bool
			Message    string
		} `graphql:"syncDeviceConfigs(input: $input)"`
	}
	err := d.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": SyncDeviceConfigsInput{SerialNumber: serial},
	})
	if err != nil {
		return err
	}
	if !mutation.SyncDeviceConfigs.Success {
		return fmt.Errorf("%d: %s", mutation.SyncDeviceConfigs.StatusCode, mutation.SyncDeviceConfigs.Message)
	}
	return nil
}

// GetDeviceBySerial retrieves a device by serial
func (d *DevicesService) GetDeviceBySerial(ctx context.Context, serial string) (*Device, error) {
	var query struct {
		Device Device `graphql:"device(serialNumber: $serial)"`
	}
	err := d.client.gql.Query(ctx, &query, map[string]interface{}{
		"serial": graphql.String(serial),
	})
	if err != nil {
		return nil, err
	}
	return &query.Device, nil
}

// GetDeviceByID retrieves a device by serial
func (d *DevicesService) GetDeviceByID(ctx context.Context, id int64) (*Device, error) {
	var query struct {
		Device Device `graphql:"device(id: $id)"`
	}
	err := d.client.gql.Query(ctx, &query, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return nil, err
	}
	return &query.Device, nil
}
