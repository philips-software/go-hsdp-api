package stl

import (
	"context"
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
