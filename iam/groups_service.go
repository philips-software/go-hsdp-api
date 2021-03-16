package iam

import (
	"fmt"
	"io"
	"net/http"
)

const (
	groupAPIVersion = "1"
)

// GetGroupOptions describes the fields on which you can search for Groups
type GetGroupOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"orgID,omitempty"`
	Name           *string `url:"name,omitempty"`
	MemberType     *string `url:"memberType,omitempty"`
}

// GroupsService implements actions on Group entities
type GroupsService struct {
	client *Client
}

// GetGroupByID retrieves a Group based on the ID
func (g *GroupsService) GetGroupByID(id string) (*Group, *Response, error) {
	req, err := g.client.newRequest(IDM, "GET", "authorize/identity/Group/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var group Group

	resp, err := g.client.do(req, &group)
	if err != nil {
		return nil, resp, err
	}
	return &group, resp, err
}

// GetGroups retrieves all groups
func (g *GroupsService) GetGroups(opt *GetGroupOptions, options ...OptionFunc) (*[]Group, *Response, error) {
	req, err := g.client.newRequest(IDM, "GET", "authorize/identity/Group", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var bundleResponse struct {
		Total int `json:"total"`
		Entry []struct {
			Resource struct {
				ID string `json:"_id"`
			} `json:"resource"`
		} `json:"entry"`
	}

	resp, err := g.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrNotFound
	}
	var groups []Group
	for _, gr := range bundleResponse.Entry {
		group, resp, err := g.GetGroupByID(gr.Resource.ID)
		if err != nil {
			return nil, resp, fmt.Errorf("GetGroups: GetGroupByID: %w", err)
		}
		groups = append(groups, *group)
	}
	return &groups, resp, nil
}

// GetGroup retrieves a Group entity based on the values passed in GetGroupOptions
func (g *GroupsService) GetGroup(opt *GetGroupOptions, options ...OptionFunc) (*Group, *Response, error) {
	req, err := g.client.newRequest(IDM, "GET", "authorize/identity/Group", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var bundleResponse struct {
		Total int `json:"total"`
		Entry []struct {
			Resource struct {
				ID string `json:"_id"`
			} `json:"resource"`
		} `json:"entry"`
	}

	resp, err := g.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrNotFound
	}
	return g.GetGroupByID(bundleResponse.Entry[0].Resource.ID)
}

// CreateGroup creates a Group
func (g *GroupsService) CreateGroup(group Group) (*Group, *Response, error) {
	if err := g.client.validate.Struct(group); err != nil {
		return nil, nil, err
	}
	req, err := g.client.newRequest(IDM, "POST", "authorize/identity/Group", &group, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var createdGroup Group

	resp, err := g.client.do(req, &createdGroup)
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
	req, err := g.client.newRequest(IDM, "PUT", "authorize/identity/Group/"+group.ID, &updateRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var updatedGroup Group

	resp, err := g.client.do(req, &updatedGroup)
	if err != nil {
		return nil, resp, err
	}
	return &group, resp, err

}

// DeleteGroup deletes the given Group
func (g *GroupsService) DeleteGroup(group Group) (bool, *Response, error) {
	req, err := g.client.newRequest(IDM, "DELETE", "authorize/identity/Group/"+group.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)

	var deleteResponse interface{}

	resp, err := g.client.do(req, &deleteResponse)
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
	req, err := g.client.newRequest(IDM, "GET", "authorize/identity/Role", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var responseStruct struct {
		Total int    `json:"total"`
		Entry []Role `json:"entry"`
	}

	resp, err := g.client.do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err
}

func (g *GroupsService) roleAction(group Group, role Role, action string) (bool, *Response, error) {
	var assignRequest = groupRequest{
		Roles: []string{role.ID},
	}
	req, err := g.client.newRequest(IDM, "POST", "authorize/identity/Group/"+group.ID+"/"+action, assignRequest, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var assignResponse interface{}

	resp, err := g.client.do(req, &assignResponse)
	if err != nil {
		return false, resp, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, resp, nil
	}
	return true, resp, err
}

// AssignRole adds a role to a group
func (g *GroupsService) AssignRole(group Group, role Role) (bool, *Response, error) {
	return g.roleAction(group, role, "$assign-role")
}

// RemoveRole removes a role from a group
func (g *GroupsService) RemoveRole(group Group, role Role) (bool, *Response, error) {
	return g.roleAction(group, role, "$remove-role")
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

type memberRequest struct {
	MemberType string   `json:"memberType"`
	Value      []string `json:"value"`
}

type groupRequest struct {
	ResourceType string      `json:"resourceType,omitempty"`
	Parameter    []Parameter `json:"parameter,omitempty"`
	Roles        []string    `json:"roles,omitempty"`
}

func (g *GroupsService) memberAction(group Group, action string, opt interface{}, options []OptionFunc) (bool, *Response, error) {
	req, err := g.client.newRequest(IDM, "POST", "authorize/identity/Group/"+group.ID+"/"+action, opt, options)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", groupAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var memberResponse interface{}

	resp, err := g.client.do(req, &memberResponse)

	if err != nil && err != io.EOF { // EOF is valid
		return false, resp, err
	}
	if resp == nil || !(resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusMultiStatus) {
		return false, resp, err
	}
	return true, resp, nil
}

func memberRequestBody(memberType string, identities ...string) memberRequest {
	var requestBody = memberRequest{
		MemberType: memberType,
		Value:      []string{},
	}
	requestBody.Value = append(requestBody.Value, identities...)
	return requestBody
}

func groupRequestBody(users ...string) groupRequest {
	var requestBody = groupRequest{
		ResourceType: "Parameters",
		Parameter: []Parameter{
			{
				Name: "UserIDCollection",
			},
		},
	}
	for _, user := range users {
		requestBody.Parameter[0].References = append(requestBody.Parameter[0].References, Reference{Reference: user})
	}
	return requestBody
}

// AddMembers adds users to the given Group
func (g *GroupsService) AddMembers(group Group, users ...string) (bool, *Response, error) {
	return g.memberAction(group, "$add-members", groupRequestBody(users...), nil)
}

// RemoveMembers removes users from the given Group
func (g *GroupsService) RemoveMembers(group Group, users ...string) (bool, *Response, error) {
	return g.memberAction(group, "$remove-members", groupRequestBody(users...), nil)
}

func addIfMatchHeader(version string) OptionFunc {
	return func(req *http.Request) error {
		req.Header.Set("If-Match", version)
		return nil
	}
}

// AddIdentities adds services to the given Group
func (g *GroupsService) AddIdentities(group Group, memberType string, identities ...string) (bool, *Response, error) {
	_, resp, err := g.GetGroupByID(group.ID)
	if err != nil {
		return false, resp, err
	}
	version := resp.Header.Get("ETag")
	return g.memberAction(group, "$assign", memberRequestBody(memberType, identities...), []OptionFunc{addIfMatchHeader(version)})
}

// RemoveIdentities removes services from the given Group
func (g *GroupsService) RemoveIdentities(group Group, memberType string, identities ...string) (bool, *Response, error) {
	_, resp, err := g.GetGroupByID(group.ID)
	if err != nil {
		return false, resp, err
	}
	version := resp.Header.Get("ETag")
	return g.memberAction(group, "$remove", memberRequestBody(memberType, identities...), []OptionFunc{addIfMatchHeader(version)})
}

// AddDevices adds services to the given Group
func (g *GroupsService) AddDevices(group Group, devices ...string) (bool, *Response, error) {
	return g.AddIdentities(group, "DEVICE", devices...)
}

// RemoveDevices removes services from the given Group
func (g *GroupsService) RemoveDevices(group Group, devices ...string) (bool, *Response, error) {
	return g.RemoveIdentities(group, "DEVICE", devices...)
}

// AddServices adds services to the given Group
func (g *GroupsService) AddServices(group Group, services ...string) (bool, *Response, error) {
	return g.AddIdentities(group, "SERVICE", services...)
}

// RemoveServices removes services from the given Group
func (g *GroupsService) RemoveServices(group Group, services ...string) (bool, *Response, error) {
	return g.RemoveIdentities(group, "SERVICE", services...)
}
