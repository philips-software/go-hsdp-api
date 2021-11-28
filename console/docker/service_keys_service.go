package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
)

type ServiceKeysService struct {
	client *Client
}

type ServiceKey struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ServiceKeyNode struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (a *ServiceKeysService) GetServiceKeys(ctx context.Context) (*[]ServiceKeyNode, error) {
	var query struct {
		Resources []struct {
			ServiceKeyNode
		} `graphql:"serviceKeys"`
	}
	err := a.client.gql.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}
	keys := make([]ServiceKeyNode, 0)
	for _, a := range query.Resources {
		keys = append(keys, a.ServiceKeyNode)
	}
	return &keys, nil
}

func (a *ServiceKeysService) GetServiceKeyByID(ctx context.Context, id int) (*ServiceKeyNode, error) {
	var query struct {
		ServiceKeyNode ServiceKeyNode `graphql:"serviceKey(id: $keyId)"`
	}
	err := a.client.gql.Query(ctx, &query, map[string]interface{}{
		"keyId": graphql.Int(id),
	})
	if err != nil {
		// Unfortunately, some regions did not receive the fix:
		// https://github.com/philips-internal/hsdp-docker-api/pull/3
		// So we fall back to the old behaviour until every region is synced up
		return a.fallbackGetServiceByKeyID(ctx, id)
	}
	return &query.ServiceKeyNode, nil
}

func (a *ServiceKeysService) CreateServiceKey(ctx context.Context, description string) (*ServiceKey, error) {
	var mutation struct {
		CreateServiceKey struct {
			ServiceKey
		} `graphql:"createServiceKey(description: $description)"`
	}
	err := a.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"description": graphql.String(description),
	})
	if err != nil {
		return nil, err
	}
	if mutation.CreateServiceKey.ServiceKey.ID == 0 {
		return nil, fmt.Errorf("error creating serviceKey")
	}
	return &mutation.CreateServiceKey.ServiceKey, nil
}

func (a *ServiceKeysService) DeleteServiceKey(ctx context.Context, key ServiceKey) error {
	var mutation struct {
		DeleteServiceKey bool `graphql:"deleteServiceKey(id: $id)"`
	}
	err := a.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"id": graphql.Int(key.ID),
	})
	if err != nil {
		return fmt.Errorf("eror deleting serviceKey: %w", err)
	}
	if !mutation.DeleteServiceKey {
		return fmt.Errorf("failed to delete serviceKey")
	}
	return nil
}

func (a *ServiceKeysService) fallbackGetServiceByKeyID(ctx context.Context, id int) (*ServiceKeyNode, error) {
	keys, err := a.GetServiceKeys(ctx)
	if err != nil {
		return nil, err
	}
	for _, k := range *keys {
		if k.ID == id {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("fallback serviceKey(id: $id) did not find a match for '%d'", id)
}
