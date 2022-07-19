package iam

import "net/http"

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

// ListSharingPoliciesOptions describes search criteria for listing RoleSharingPolicy resources
type ListSharingPoliciesOptions struct {
	TargetOrganizationID *string `url:"targetOrganizationId,omitempty"`
	SharingPolicy        *string `url:"sharingPolicy,omitempty"`
	RecordsPerPage       *int    `url:"recordsPerPage,omitempty"`
	StartPage            *int    `url:"startPage,omitempty"`
}

// RoleSharingPolicy describes a role sharing policy
type RoleSharingPolicy struct {
	SharingPolicy        string `json:"sharingPolicy"`
	Purpose              string `json:"purpose"`
	TargetOrganizationID string `json:"targetOrganizationId"`
	InternalID           string `json:"internalId,omitempty"`
	SourceOrganizationID string `json:"sourceOrganizationId,omitempty"`
	Meta                 *Meta  `json:"meta,omitempty"`
}

// GetRoles retries based on GetRolesOptions
func (p *RolesService) GetRoles(opt *GetRolesOptions) (*[]Role, *Response, error) {
	req, err := p.client.newRequest(IDM, http.MethodGet, "authorize/identity/Role", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var responseStruct struct {
		Total int    `json:"total"`
		Entry []Role `json:"entry"`
	}

	resp, err := p.client.do(req, &responseStruct)
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
	req, err := p.client.newRequest(IDM, http.MethodGet, "authorize/identity/Role/"+roleID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var role Role

	resp, err := p.client.do(req, &role)
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
	req, _ := p.client.newRequest(IDM, http.MethodPost, "authorize/identity/Role", role, nil)
	req.Header.Set("api-version", roleAPIVersion)

	var createdRole Role

	resp, err := p.client.do(req, &createdRole)
	if err != nil {
		return nil, resp, err
	}
	return &createdRole, resp, nil
}

type RoleResponse map[string]interface{}

// DeleteRole deletes the given Role
func (p *RolesService) DeleteRole(role Role) (RoleResponse, *Response, error) {
	req, err := p.client.newRequest(IDM, http.MethodDelete, "authorize/identity/Role/"+role.ID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var roleResponse RoleResponse

	resp, err := p.client.do(req, &roleResponse)

	return roleResponse, resp, err
}

// GetRolePermissions retrieves the permissions associated with the Role
func (p *RolesService) GetRolePermissions(role Role) (*[]string, *Response, error) {
	opt := &GetRolesOptions{RoleID: &role.ID}

	req, err := p.client.newRequest(IDM, http.MethodGet, "authorize/identity/Permission", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var permissionResponse struct {
		Total int          `json:"total"`
		Entry []Permission `json:"entry"`
	}

	resp, err := p.client.do(req, &permissionResponse)
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
func (p *RolesService) rolePermissionAction(role Role, permissions []string, action string) (RoleResponse, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = permissions

	req, err := p.client.newRequest(IDM, http.MethodPost, "authorize/identity/Role/"+role.ID+"/"+action, &permissionRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var roleResponse RoleResponse

	resp, err := p.client.do(req, &roleResponse)
	if err != nil {
		return roleResponse, resp, err
	}
	return roleResponse, resp, nil

}

func (p *RolesService) AddRolePermission(role Role, permission string) (RoleResponse, *Response, error) {
	return p.rolePermissionAction(role, []string{permission}, "$assign-permission")
}

// RemoveRolePermission removes the permission from the Role
func (p *RolesService) RemoveRolePermission(role Role, permission string) (RoleResponse, *Response, error) {
	return p.rolePermissionAction(role, []string{permission}, "$remove-permission")
}

func (p *RolesService) ApplySharingPolicy(role Role, policy RoleSharingPolicy) (*RoleSharingPolicy, *Response, error) {
	req, err := p.client.newRequest(IDM, http.MethodPut, "authorize/identity/Role/"+role.ID+"/"+"$apply-sharing-policy", &policy, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var roleResponse RoleSharingPolicy

	resp, err := p.client.do(req, &roleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &roleResponse, resp, nil
}

func (p *RolesService) RemoveSharingPolicy(role Role, policy RoleSharingPolicy) (*RoleSharingPolicy, *Response, error) {
	req, err := p.client.newRequest(IDM, http.MethodPost, "authorize/identity/Role/"+role.ID+"/"+"$remove-sharing-policy", &policy, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	var roleResponse RoleSharingPolicy

	resp, err := p.client.do(req, &roleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &roleResponse, resp, nil
}

func (p *RolesService) ListSharingPolicies(role Role, opt *ListSharingPoliciesOptions) (*[]RoleSharingPolicy, *Response, error) {
	var listResponse struct {
		Count  int                 `json:"count"`
		Result []RoleSharingPolicy `json:"result"`
	}

	req, err := p.client.newRequest(IDM, http.MethodGet, "authorize/identity/Role/"+role.ID+"/"+"$list-sharing-policies", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", roleAPIVersion)

	resp, err := p.client.do(req, &listResponse)
	if err != nil {
		return nil, resp, err
	}
	return &listResponse.Result, resp, nil
}
