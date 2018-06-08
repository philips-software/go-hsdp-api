package iam

import (
	"net/http"
)

const (
	groupAPIVersion = "1"
)

// GetGroupOptions describes the fileds on which you can search for Groups
type GetGroupOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"Id,omitempty"`
	Name           *string `url:"name,omitempty"`
}

// GroupsService implements actions on Group entities
type GroupsService struct {
	client *Client
}

// GetGroupByID retrieves a Group based on the ID
func (g *GroupsService) GetGroupByID(id string) (*Group, *Response, error) {
	return g.GetGroup(&GetGroupOptions{ID: &id}, nil)
}

// GetGroup retrieves a Group entity based on the values passed in GetGroupOptions
func (g *GroupsService) GetGroup(opt *GetGroupOptions, options ...OptionFunc) (*Group, *Response, error) {
	req, err := g.client.NewIDMRequest("GET", "authorize/identity/Group", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var bundleResponse interface{}

	resp, err := g.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var group Group
	err = group.parseFromBundle(bundleResponse)
	return &group, resp, err
}

// CreateGroup creates a Group
func (g *GroupsService) CreateGroup(group Group) (*Group, *Response, error) {
	if err := group.Validate(); err != nil {
		return nil, nil, err
	}
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group", &group, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var createdGroup Group

	resp, err := g.client.Do(req, &createdGroup)
	if err != nil {
		return nil, resp, err
	}
	return &createdGroup, resp, err

}

// UpdateGroup updates the Group
func (g *GroupsService) UpdateGroup(group Group) (*Group, *Response, error) {
	var updateRequest struct {
		Description string `json:"description"`
	}
	updateRequest.Description = group.Description
	req, err := g.client.NewIDMRequest("PUT", "authorize/identity/Group/"+group.ID, &updateRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var updatedGroup Group

	resp, err := g.client.Do(req, &updatedGroup)
	if err != nil {
		return nil, resp, err
	}
	return &group, resp, err

}

// DeleteGroup deletes the given Group
func (g *GroupsService) DeleteGroup(group Group) (bool, *Response, error) {
	req, err := g.client.NewIDMRequest("DELETE", "authorize/identity/Group/"+group.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var deleteResponse interface{}

	resp, err := g.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err

}

func (g *GroupsService) AssignRole(group Group, role Role) (bool, *Response, error) {
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group/"+group.ID+"/$assign-role", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	var assignRequest struct {
		roles []string `json:"roles"`
	}
	assignRequest.roles = []string{role.ID}

	var assignResponse interface{}

	resp, err := g.client.Do(req, &assignResponse)
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, nil
	}
	return true, resp, err
}

func (g *GroupsService) RemoveRole(group Group, role Role) (bool, *Response, error) {
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group/"+group.ID+"/$remove-role", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	var removeRequest struct {
		roles []string `json:"roles"`
	}
	removeRequest.roles = []string{role.ID}

	var removeResponse interface{}

	resp, err := g.client.Do(req, &removeResponse)
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, nil
	}
	return true, resp, err
}
