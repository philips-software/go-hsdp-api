package iam

type RolesService struct {
	client *Client
}

type GetRolesOptions struct {
	Name           *string `url:"name,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	RoleID         *string `url:"roleId,omitempty"`
}

func (p *RolesService) GetRolesByGroupID(groupID string) (*[]Role, *Response, error) {
	opt := &GetRolesOptions{
		GroupID: &groupID,
	}
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Role", opt, nil)
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

func (p *RolesService) GetRoleByID(roleID string) (*Role, *Response, error) {
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Role/"+roleID, nil, nil)
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
	return &role, resp, err
}

// TODO: below method actually is not yet implemented on the HSDP side
func (p *RolesService) UpdateRole(role *Role) (*Role, *Response, error) {
	var updateRoleRequestBody struct {
		Description string `json:"description"`
	}
	updateRoleRequestBody.Description = role.Description
	req, err := p.client.NewIDMRequest("PUT", "authorize/identity/Role", &updateRoleRequestBody, nil)
	req.Header.Set("api-version", "1")

	var updatedRole Role
	resp, err := p.client.Do(req, &updatedRole)

	if err != nil {
		return nil, resp, err
	}
	return role, resp, nil

}

func (p *RolesService) CreateRole(name, description, managingOrganization string) (*Role, *Response, error) {
	role := &Role{
		Name:                 name,
		Description:          description,
		ManagingOrganization: managingOrganization,
	}
	req, err := p.client.NewIDMRequest("POST", "authorize/identity/Role", role, nil)
	req.Header.Set("api-version", "1")

	var createdRole Role

	resp, err := p.client.Do(req, &createdRole)
	if err != nil {
		return nil, resp, err
	}
	return &createdRole, resp, err
}

func (p *RolesService) GetRole(opt *GetRolesOptions, options ...OptionFunc) (*Role, *Response, error) {
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Role", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse interface{}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var role Role
	err = role.parseFromBundle(bundleResponse)
	return &role, resp, err
}

func (p *RolesService) GetRolePermissions(role Role) (*[]string, error) {
	opt := &GetRolesOptions{RoleID: &role.ID}

	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Permission", opt, nil)
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

func (p *RolesService) AddRolePermission(role Role, permission string) (*Role, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = []string{permission}

	req, err := p.client.NewIDMRequest("POST", "authorize/identity/Role/"+role.ID+"/$assign-permission", &permissionRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse interface{}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return nil, resp, err

}

func (p *RolesService) RemoveRolePermission(role Role, permission string) (*Role, *Response, error) {
	var permissionRequest struct {
		Permissions []string `json:"permissions"`
	}
	permissionRequest.Permissions = []string{permission}

	req, err := p.client.NewIDMRequest("POST", "authorize/identity/Role/"+role.ID+"/$remove-permission", &permissionRequest, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse interface{}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return nil, resp, err
}
