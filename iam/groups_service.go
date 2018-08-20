package iam

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jeffail/gabs"
)

const (
	groupAPIVersion = "1"
)

// Group represents an IAM group resource
type Group struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	ManagingOrganization string `json:"managingOrganization,omitempty"`
}

func (g *Group) validate() error {
	if g.ManagingOrganization == "" {
		return errMissingManagingOrganization
	}
	if g.Name == "" {
		return errMissingName
	}
	return nil
}

func (g *Group) parseFromBundle(v interface{}) error {
	m, _ := json.Marshal(v)
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return err
	}
	r := jsonParsed.Path("entry").Index(0).Path("resource")
	g.ID = r.Path("_id").Data().(string)
	g.ManagingOrganization, _ = r.Path("orgId").Data().(string)
	g.Name, _ = r.Path("groupName").Data().(string)
	g.Description, _ = r.Path("groupDescription").Data().(string)
	return nil
}

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

	resp, err := g.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var group Group
	err = group.parseFromBundle(bundleResponse)
	return &group, resp, err
}

// CreateGroup creates a Group
func (g *GroupsService) CreateGroup(group Group) (*Group, *Response, error) {
	if err := group.validate(); err != nil {
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
		return false, resp, err
	}
	return true, resp, nil

}

// GetRoles returns the roles assigned to this group
func (g *GroupsService) GetRoles(group Group) (*[]Role, *Response, error) {
	opt := &GetRolesOptions{
		GroupID: &group.ID,
	}
	req, err := g.client.NewIDMRequest("GET", "authorize/identity/Role", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var responseStruct struct {
		Total int    `json:"total"`
		Entry []Role `json:"entry"`
	}

	resp, err := g.client.Do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err
}

// AssignRole adds a role to a group
func (g *GroupsService) AssignRole(group Group, role Role) (bool, *Response, error) {
	var assignRequest struct {
		Roles []string `json:"roles"`
	}
	assignRequest.Roles = []string{role.ID}
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group/"+group.ID+"/$assign-role", assignRequest, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var assignResponse interface{}

	resp, err := g.client.Do(req, &assignResponse)
	if err != nil {
		return false, resp, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, nil
	}
	return true, resp, err
}

// RemoveRole removes a role from a group
func (g *GroupsService) RemoveRole(group Group, role Role) (bool, *Response, error) {
	var removeRequest struct {
		Roles []string `json:"roles"`
	}
	removeRequest.Roles = []string{role.ID}
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group/"+group.ID+"/$remove-role", removeRequest, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var removeResponse interface{}

	resp, err := g.client.Do(req, &removeResponse)
	if err != nil {
		return false, resp, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, nil
	}
	return true, resp, err
}

// Reference holds a reference
type Reference struct {
	Reference string `json:"reference"`
}

// Parameter holds named references
type Parameter struct {
	Name       string      `json:"name"`
	References []Reference `json:"references"`
}

// AddUser adds a user to the given Group
func (g *GroupsService) AddUser(group Group, userID string) (bool, *Response, error) {

	var addRequest = struct {
		ResourceType string      `json:"resourceType"`
		Parameter    []Parameter `json:"parameter"`
	}{
		ResourceType: "Parameters",
		Parameter: []Parameter{
			{
				Name: "UserIDCollection",
				References: []Reference{
					{Reference: userID},
				},
			},
		},
	}
	req, err := g.client.NewIDMRequest("POST", "authorize/identity/Group/"+group.ID+"/$add-members", addRequest, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var addResponse interface{}

	resp, err := g.client.Do(req, &addResponse)

	if err != nil && err != io.EOF { // EOF is valid
		return false, resp, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, err
	}
	return true, resp, nil
}
