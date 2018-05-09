package iam

import (
	"net/http"
)

const (
	GroupAPIVersion = "1"
)

type GetGroupOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"Id,omitempty"`
	Name           *string `url:"name,omitempty"`
}

type GroupsService struct {
	client *Client
}

func (g *GroupsService) GetGroupByID(id string) (*Group, *Response, error) {
	return g.GetGroup(&GetGroupOptions{ID: &id}, nil)
}

func (g *GroupsService) GetGroup(opt *GetGroupOptions, options ...OptionFunc) (*Group, *Response, error) {
	req, err := g.client.NewIDMRequest("GET", "authorize/identity/Group", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", GroupAPIVersion)

	var bundleResponse interface{}

	resp, err := g.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var group Group
	group.ParseFromBundle(bundleResponse)
	return &group, resp, err
}

func (g *GroupsService) CreateGroup(group Group) (*Group, *Response, error) {
	if err := group.Validate(); err != nil {
		return nil, nil, err
	}
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group", &group, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", GroupAPIVersion)

	var createdGroup Group

	resp, err := g.client.Do(req, &createdGroup)
	if err != nil {
		return nil, resp, err
	}
	return &createdGroup, resp, err

}

func (g *GroupsService) UpdateGroup(group Group) (*Group, *Response, error) {
	var updateRequest struct {
		Description string `json:"description"`
	}
	updateRequest.Description = group.Description
	req, err := g.client.NewIDMRequest("PUT", "authorize/identity/Group/"+group.ID, &updateRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", GroupAPIVersion)

	var updatedGroup Group

	resp, err := g.client.Do(req, &updatedGroup)
	if err != nil {
		return nil, resp, err
	}
	return &group, resp, err

}

func (g *GroupsService) DeleteGroup(group Group) (bool, *Response, error) {
	req, err := g.client.NewIDMRequest("DELETE", "authorize/identity/Group/"+group.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", GroupAPIVersion)

	var deleteResponse interface{}

	resp, err := g.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err

}
