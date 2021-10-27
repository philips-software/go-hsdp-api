package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
)

type RepositoriesService struct {
	client *Client
}

type Repository struct {
	ID           string                 `json:"id"`
	NamespaceId  string                 `json:"namespaceId"`
	Name         string                 `json:"name"`
	NumPulls     int                    `json:"numPulls"`
	NumTags      int                    `json:"numTags"`
	LastPushedAt time.Time              `json:"lastPushedAt"`
	Details      RepositoryDetailsInput `json:"details"`
}

type Tag struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	CompressedSize int       `json:"compressedSize"`
	UpdatedAt      time.Time `json:"updatedAt"`
	NumPulls       int       `json:"numPulls"`
	Digest         string    `json:"digest"`
	ImageId        string    `json:"imageId"`
}

type RepositoryInput struct {
	NamespaceID string `json:"namespaceId"`
	Name        string `json:"name"`
}

type RepositoryDetailsInput struct {
	ShortDescription string `json:"shortDescription"`
	FullDescription  string `json:"fullDescription"`
}

type RepositoryResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NamespaceID string `json:"namespaceId"`
}

func (r *RepositoriesService) CreateRepository(ctx context.Context, repository RepositoryInput, details RepositoryDetailsInput) (*RepositoryResult, error) {
	var mutation struct {
		Repository RepositoryResult `graphql:"createRepository(repository: $repository, details: $details)"`
	}
	err := r.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"repository": repository,
		"details":    details,
	})
	if err != nil {
		return nil, err
	}
	if mutation.Repository.ID == "" {
		return nil, fmt.Errorf("error creating repository")
	}
	return &mutation.Repository, nil
}

func (r *RepositoriesService) UpdateRepository(ctx context.Context, repository Repository, details RepositoryDetailsInput) (*RepositoryDetailsInput, error) {
	var mutation struct {
		Resources struct {
			Details RepositoryDetailsInput
		} `graphql:"updateRepository(id: $id, details: $details)"`
	}
	err := r.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"id":      graphql.String(repository.ID),
		"details": details,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.Resources.Details, nil
}

func (r *RepositoriesService) DeleteRepository(ctx context.Context, repository Repository) error {
	var mutation struct {
		DeleteRepository bool `graphql:"deleteRepository(id: $id)"`
	}
	err := r.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"id": graphql.String(repository.ID),
	})
	if err != nil {
		return fmt.Errorf("eror deleting repository: %w", err)
	}
	if !mutation.DeleteRepository {
		return fmt.Errorf("failed to delete repository")
	}
	return nil
}

func (r *RepositoriesService) GetRepository(ctx context.Context, namespaceId, name string) (*Repository, error) {
	var query struct {
		Repository Repository `graphql:"repository(namespaceId: $namespaceId, name: $name)"`
	}
	err := r.client.gql.Query(ctx, &query, map[string]interface{}{
		"namespaceId": graphql.String(namespaceId),
		"name":        graphql.String(name),
	})
	if err != nil {
		return nil, err
	}
	return &query.Repository, nil
}

func (r *RepositoriesService) GetTags(ctx context.Context, repositoryId string) (*[]Tag, error) {
	var query struct {
		Tags []Tag `graphql:"tags(repositoryId: $repositoryId, page: $page, limit: $limit, orderBy: UPDATED_AT)"`
	}
	err := r.client.gql.Query(ctx, &query, map[string]interface{}{
		"repositoryId": graphql.String(repositoryId),
		"page":         graphql.Int(1),
		"limit":        graphql.Int(1000),
	})
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0)
	tags = append(tags, query.Tags...)
	return &tags, nil
}
