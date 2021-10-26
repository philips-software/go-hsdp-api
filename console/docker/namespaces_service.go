package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
)

type NamespacesService struct {
	client *Client
}

type Namespace struct {
	ID           string                   `json:"id"`
	IsPublic     bool                     `json:"isPublic"`
	CreatedAt    time.Time                `json:"createdAt"`
	IsOrgSpaceID bool                     `json:"isOrgSpaceId"`
	NumRepos     int                      `json:"numRepos"`
	UserAccess   UserNamespaceAccessInput `json:"userAccess"`
}

type UserNamespaceAccessInput struct {
	CanPull   bool `json:"canPull"`
	CanPush   bool `json:"canPush"`
	CanDelete bool `json:"canDelete"`
	IsAdmin   bool `json:"isAdmin"`
}

type UpdateUserNamespacesAccessInput struct {
	ID int `json:"id"`
	UserNamespaceAccessInput
}

type NamespaceAccess struct {
	ID int `json:"id"`
	UserNamespaceAccessInput
}

type NamespaceUser struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Username        string          `json:"username"`
	NamespaceAccess NamespaceAccess `json:"namespaceAccess"`
}

type NamespaceUserResult struct {
	ID     int    `json:"id"`
	UserID string `json:"userId"`
	UserNamespaceAccessInput
}

func (s *NamespacesService) GetNamespaces(ctx context.Context) (*[]Namespace, error) {
	var query struct {
		Resources []struct {
			Namespace
		} `graphql:"namespaces(userId: $userId, page: $page, limit: $limit)"`
	}
	userID, err := s.client.UserID()
	if err != nil {
		return nil, fmt.Errorf("userId error: %w", err)
	}
	err = s.client.gql.Query(ctx, &query, map[string]interface{}{
		"userId": graphql.String(userID),
		"page":   graphql.Int(1),
		"limit":  graphql.Int(1000),
	})
	if err != nil {
		return nil, err
	}
	namespaces := make([]Namespace, 0)
	for _, a := range query.Resources {
		namespaces = append(namespaces, a.Namespace)
	}
	return &namespaces, nil
}

type NamespaceInput struct {
	ID graphql.String `json:"id"`
}

func (s *NamespacesService) CreateNamespace(ctx context.Context, id string) (*Namespace, error) {
	var mutation struct {
		CreateNamespace NamespaceInput `graphql:"createNamespace(namespace: $namespace)"`
	}
	err := s.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"namespace": NamespaceInput{
			ID: graphql.String(id),
		},
	})
	if err != nil {
		return nil, err
	}
	if mutation.CreateNamespace.ID == "" {
		return nil, fmt.Errorf("error creating namespace")
	}
	return s.GetNamespaceByID(ctx, id)
}

func (s *NamespacesService) GetNamespaceByID(ctx context.Context, id string) (*Namespace, error) {
	var query struct {
		Namespace Namespace `graphql:"namespace(id: $namespaceId)"`
	}
	err := s.client.gql.Query(ctx, &query, map[string]interface{}{
		"namespaceId": graphql.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("namespace read: %w", err)
	}
	return &query.Namespace, nil
}

func (s *NamespacesService) GetNamespaceUsers(ctx context.Context, ns Namespace) (*[]NamespaceUser, error) {
	var query struct {
		Resources []struct {
			NamespaceUser
		} `graphql:"namespaceUsers(namespaceId: $namespaceId)"`
	}
	err := s.client.gql.Query(ctx, &query, map[string]interface{}{
		"namespaceId": graphql.String(ns.ID),
	})
	if err != nil {
		return nil, err
	}
	namespaceUsers := make([]NamespaceUser, 0)
	for _, a := range query.Resources {
		namespaceUsers = append(namespaceUsers, a.NamespaceUser)
	}
	return &namespaceUsers, nil
}

func (s *NamespacesService) AddNamespaceUser(ctx context.Context, namespaceID, username string, access UserNamespaceAccessInput) (*NamespaceUserResult, error) {
	var mutation struct {
		AddUserToNamespace NamespaceUserResult `graphql:"addUserToNamespace(namespaceId: $namespaceId, username: $username, userAccess: $access)"`
	}
	err := s.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"namespaceId": graphql.String(namespaceID),
		"username":    graphql.String(username),
		"access":      access,
	})
	if err != nil {
		return nil, err
	}
	if mutation.AddUserToNamespace.ID == 0 {
		return nil, fmt.Errorf("error add user to namespace")
	}
	return &mutation.AddUserToNamespace, nil
}

func (s *NamespacesService) GetRepositories(ctx context.Context, namespaceId string) (*[]Repository, error) {
	var query struct {
		Resources struct {
			Repositories []Repository
		} `graphql:"namespace(id: $namespaceId)"`
	}
	err := s.client.gql.Query(ctx, &query, map[string]interface{}{
		"namespaceId": graphql.String(namespaceId),
	})
	if err != nil {
		return nil, err
	}
	repositories := make([]Repository, 0)
	repositories = append(repositories, query.Resources.Repositories...)
	return &repositories, nil
}

func (s *NamespacesService) UpdateNamespaceUserAccess(ctx context.Context, id int, access UserNamespaceAccessInput) error {
	var mutation struct {
		UpdateNamespaceUserAccess NamespaceUserResult `graphql:"updateUserNamespaceAccess(id: $id, userAccess: $access)"`
	}
	err := s.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"id":     graphql.Int(id),
		"access": access,
	})
	if err != nil {
		return err
	}
	if mutation.UpdateNamespaceUserAccess.ID == 0 {
		return fmt.Errorf("error updating namespace user access")
	}
	return nil
}

func (s *NamespacesService) DeleteNamespaceUser(ctx context.Context, namespaceID, userId string) error {
	var mutation struct {
		DeleteUserFromNamespace bool `graphql:"removeUserFromNamespace(namespaceId: $namespaceId, userId: $userId)"`
	}
	err := s.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"namespaceId": graphql.String(namespaceID),
		"userId":      graphql.String(userId),
	})
	if err != nil {
		return fmt.Errorf("eror removing user from namespace: %w", err)
	}
	if !mutation.DeleteUserFromNamespace {
		return fmt.Errorf("failed to remove user from namespace")
	}
	return nil
}

func (s *NamespacesService) DeleteNamespace(ctx context.Context, ns Namespace) error {
	var mutation struct {
		DeleteNamespace bool `graphql:"deleteNamespace(id: $id)"`
	}
	err := s.client.gql.Mutate(ctx, &mutation, map[string]interface{}{
		"id": graphql.String(ns.ID),
	})
	if err != nil {
		return fmt.Errorf("eror deleting namespace: %w", err)
	}
	if !mutation.DeleteNamespace {
		return fmt.Errorf("failed to delete namespace")
	}
	return nil
}
