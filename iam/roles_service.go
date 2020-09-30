package iam

import (
	"bytes"
	"net/http"
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

// GetRoles retries based on GetRolesOptions
func (p *RolesService) GetRoles(opt *GetRolesOptions) (*[]Role, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Role", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
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

// GetRolesByGroupID retrieves Roles based on group ID
func (p *RolesService) GetRolesByGroupID(groupID string) (*[]Role, *Response, error) {
	opt := &GetRolesOptions{
		GroupID: &groupID,
	}
	return p.GetRoles(opt)
}

// GetRoleByID retrieves a role by ID
func (p *RolesService) GetRoleByID(roleID string) (*Role, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Role/"+roleID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var role Role

	resp, err := p.client.Do(req, &role)
	if err != nil {
		return nil, resp, err
	}
	if role.ID != roleID {
		return nil, resp, ErrNotFound
	}
	return &role, resp, err
}

// CreateRole creates a Role
func (p *RolesService) CreateRole(name, description, managingOrganization string) (*Role, *Response, error) {
	role := &Role{
		Name:                 name,
		Description:          description,
		ManagingOrganization: managingOrganization,
	}
	req, _ := p.client.NewRequest(IDM, "POST", "authorize/identity/Role", role, nil)
	req.Header.Set("api-version", roleAPIVersion)

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
	req.Header.Set("api-version", roleAPIVersion)

	var deleteResponse bytes.Buffer

	resp, err := p.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

// GetRolePermissions retrieves the permissions associated with the Role
func (p *RolesService) GetRolePermissions(role Role) (*[]string, *Response, error) {
	opt := &GetRolesOptions{RoleID: &role.ID}

	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Permission", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var permissionResponse struct {
		Total int          `json:"total"`
		Entry []Permission `json:"entry"`
	}

	resp, err := p.client.Do(req, &permissionResponse)
	if err != nil {
		return nil, resp, err
	}
	var permissions []string
	for _, p := range permissionResponse.Entry {
		permissions = append(permissions, p.Name)
	}
	return &permissions, resp, err

}

// AddRolePermission adds a given permission to the Role
func (p *RolesService) rolePermissionAction(role Role, permission string, action string) (bool, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = []string{permission}

	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Role/"+role.ID+"/"+action, &permissionRequest, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	return true, resp, err

}

func (p *RolesService) AddRolePermission(role Role, permission string) (bool, *Response, error) {
	return p.rolePermissionAction(role, permission, "$assign-permission")
}

// RemoveRolePermission removes the permission from the Role
func (p *RolesService) RemoveRolePermission(role Role, permission string) (bool, *Response, error) {
	return p.rolePermissionAction(role, permission, "$remove-permission")
}
