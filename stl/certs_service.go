package stl

import (
	"context"
	"github.com/hasura/go-graphql-client"
)

type CertsService struct {
	client *Client
}

type CustomCert struct {
	ID       int64  `json:"id"`
	DeviceID int64  `json:"deviceId,omitempty"`
	Name     string `json:"name"`
	Key      string `json:"key"`
	Cert     string `json:"cert"`
}

type CreateAppCustomCertInput struct {
	CustomCert
	SerialNumber string `json:"serialNumber"`
}

type UpdateAppCustomCertInput struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Cert string `json:"cert"`
}

type DeleteAppCustomCertInput struct {
	ID int64 `json:"id"`
}

func (a *CertsService) GetCustomCertByID(ctx context.Context, id int64) (*CustomCert, error) {
	var query struct {
		CustomCert CustomCert `graphql:"appCustomCert(id: $id)"`
	}
	err := a.client.gql.Query(ctx, &query, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return nil, err
	}
	return &query.CustomCert, nil
}

func (a *CertsService) GetCustomCertsBySerial(ctx context.Context, serial string) (*[]CustomCert, error) {
	var query struct {
		Resources struct {
			Edges []struct {
				Node CustomCert
			}
		} `graphql:"appCustomCerts(serialNumber: $serial, first: 10000)"`
	}
	err := a.client.gql.Query(ctx, &query, map[string]interface{}{
		"serial": graphql.String(serial),
	})
	if err != nil {
		return nil, err
	}
	certs := make([]CustomCert, 0)
	for _, a := range query.Resources.Edges {
		certs = append(certs, a.Node)
	}
	return &certs, nil
}

func (a *CertsService) CreateCustomCert(ctx context.Context, input CreateAppCustomCertInput) (*CustomCert, error) {
	var mutation struct {
		CreateAppCustomCert struct {
			Success    bool
			Message    string
			StatusCode int
			RequestID  string
			CustomCert CustomCert
		} `graphql:"createAppCustomCert(input: $input)"`
	}
	err := a.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.CreateAppCustomCert.CustomCert, nil
}

func (a *CertsService) UpdateCustomCert(ctx context.Context, input UpdateAppCustomCertInput) (*CustomCert, error) {
	var mutation struct {
		UpdateApplicationResource struct {
			Success    bool
			Message    string
			StatusCode int
			RequestID  string
			CustomCert CustomCert
		} `graphql:"updateAppCustomCert(input: $input)"`
	}
	err := a.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.UpdateApplicationResource.CustomCert, nil
}

func (a *AppsService) DeleteCustomCert(ctx context.Context, input DeleteAppCustomCertInput) (bool, error) {
	var mutation struct {
		DeleteAppCustomCert struct {
			Success    bool
			Message    string
			StatusCode int
		} `graphql:"deleteAppCustomCert(input: $input)"`
	}
	err := a.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
