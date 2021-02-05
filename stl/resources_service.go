package stl

import (
	"context"
	"github.com/hasura/go-graphql-client"
)

type ResourcesService struct {
	client *Client
}

type ApplicationResourceType string

type Resource struct {
}

func (r *ResourcesService) GetResourcesBySerialAndType(ctx context.Context, serial, resourceType string) (*[]Resource, error) {
	var query struct {
		Resources []Resource `graphql:"resources(serialNumber: $serial, type: $type)"`
	}
	err := r.client.gql.Query(ctx, &query, map[string]interface{}{
		"serial": graphql.String(serial),
		"type":   ApplicationResourceType(resourceType),
	})
	if err != nil {
		return nil, err
	}
	return &query.Resources, nil
}
