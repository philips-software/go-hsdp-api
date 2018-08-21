package iam

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/jeffail/gabs"
)

var (
	roleAPIVersion = "1"
)

// Role represents an IAM resource
type Role struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	ManagingOrganization string `json:"managingOrganization"`
}

// RolesService provides operations on IAM roles resources
type RolesService struct {
	client *Client
}

// GetRolesOptions describes search criteria for looking up roles
type GetRolesOptions struct {
	Name           *string `url:"name,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	RoleID         *string `url:"roleId,omitempty"`
}

// GetRolesByGroupID retrieves Roles based on group ID
func (p *RolesService) GetRolesByGroupID(groupID string) (*[]Role, *Response, error) {
	opt := &GetRolesOptions{
		GroupID: &groupID,
	}
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Role", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var responseStruct struct {
		Total int    `json:"total"`
		Entry []Role `json:"entry"`
	}

	resp, err := p.client.Do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err

}

// GetRoleByID retrieves a role by ID
func (p *RolesService) GetRoleByID(roleID string) (*Role, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Role/"+roleID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")
	req.Header.Set("Content-Type", "application/json")

	var role Role

	resp, err := p.client.Do(req, &role)
	if err != nil {
		return nil, resp, err
	}
	if role.ID != roleID {
		return nil, resp, errNotFound
	}
	return &role, resp, err
}

// UpdateRole updates a role
// TODO: below method actually is not yet implemented on the HSDP side
func (p *RolesService) UpdateRole(role *Role) (*Role, *Response, error) {
	var updateRoleRequestBody struct {
		Description string `json:"description"`
	}
	updateRoleRequestBody.Description = role.Description
	req, err := p.client.NewRequest(IDM, "PUT", "authorize/identity/Role", &updateRoleRequestBody, nil)
	req.Header.Set("api-version", "1")

	var updatedRole Role
	resp, err := p.client.Do(req, &updatedRole)

	if err != nil {
		return nil, resp, err
	}
	return role, resp, nil

}

// CreateRole creates a Role
func (p *RolesService) CreateRole(name, description, managingOrganization string) (*Role, *Response, error) {
	role := &Role{
		Name:                 name,
		Description:          description,
		ManagingOrganization: managingOrganization,
	}
	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Role", role, nil)
	req.Header.Set("api-version", "1")

	var createdRole Role

	resp, err := p.client.Do(req, &createdRole)
	if err != nil {
		return nil, resp, err
	}
	return &createdRole, resp, err
}

// DeleteRole deletes the given Role
func (p *RolesService) DeleteRole(role Role) (bool, *Response, error) {
	req, err := p.client.NewRequest(IDM, "DELETE", "authorize/identity/Role/"+role.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", "1")

	var deleteResponse interface{}

	resp, err := p.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

// GetRole retrieve a Role resource based on GetRoleOptions parameters
func (p *RolesService) GetRole(opt *GetRolesOptions, options ...OptionFunc) (*Role, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Role", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	roles, err := p.parseFromBundle(bundleResponse.Bytes())
	if err != nil {
		return nil, resp, err
	}
	return &(*roles)[0], resp, err
}

// GetRolePermissions retrieves the permissions assosciates with the Role
func (p *RolesService) GetRolePermissions(role Role) (*[]string, error) {
	opt := &GetRolesOptions{RoleID: &role.ID}

	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Permission", opt, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("api-version", "1")

	var permissionResponse struct {
		Total int    `json:"total"`
		Entry []Role `json:"entry"`
	}

	_, err = p.client.Do(req, &permissionResponse)
	if err != nil {
		return nil, err
	}
	var permissions []string
	for _, p := range permissionResponse.Entry {
		permissions = append(permissions, p.Name)
	}
	return &permissions, err

}

// AddRolePermission adds a given permission to the Role
func (p *RolesService) AddRolePermission(role Role, permission string) (*Role, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = []string{permission}

	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Role/"+role.ID+"/$assign-permission", &permissionRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return nil, resp, err

}

// RemoveRolePermission removes the permission from the Role
func (p *RolesService) RemoveRolePermission(role Role, permission string) (*Role, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = []string{permission}

	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Role/"+role.ID+"/$remove-permission", &permissionRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return nil, resp, err
}

func (p *RolesService) parseFromBundle(bundle []byte) (*[]Role, error) {
	jsonParsed, err := gabs.ParseJSON(bundle)
	if err != nil {
		return nil, err
	}
	count, ok := jsonParsed.S("total").Data().(float64)
	if !ok || count == 0 {
		return nil, errors.New("empty result")
	}
	roles := make([]Role, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var p Role
		p.ID = r.Path("id").Data().(string)
		p.ManagingOrganization, _ = r.Path("managingOrganization").Data().(string)
		p.Name, _ = r.Path("name").Data().(string)
		p.Description, _ = r.Path("description").Data().(string)
		roles[i] = p
	}
	return &roles, nil
}
