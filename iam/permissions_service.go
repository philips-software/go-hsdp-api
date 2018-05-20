package iam

const permissionAPIVersion = "1"

type PermissionsService struct {
	client *Client
}

type GetPermissionOptions struct {
	ID     *string `url:"_id,omitempty"`
	Name   *string `url:"name,omitempty"`
	RoleID *string `url:"roleId,omitempty"`
}

func (p *PermissionsService) GetPermissionByID(id string) (*Permission, *Response, error) {
	return p.GetPermission(&GetPermissionOptions{ID: &id}, nil)
}

func (p *PermissionsService) GetPermissionByName(name string) (*Permission, *Response, error) {
	return p.GetPermission(&GetPermissionOptions{Name: &name}, nil)
}

func (p *PermissionsService) GetPermissionsByRoleID(roleID string) (*[]Permission, *Response, error) {
	opt := &GetPermissionOptions{
		RoleID: &roleID,
	}
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Permission", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", permissionAPIVersion)

	var responseStruct struct {
		Total int          `json:"total"`
		Entry []Permission `json:"entry"`
	}

	resp, err := p.client.DoSigned(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err
}

func (p *PermissionsService) GetPermission(opt *GetPermissionOptions, options ...OptionFunc) (*Permission, *Response, error) {
	req, err := p.client.NewIDMRequest("GET", "authorize/identity/Permission", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", permissionAPIVersion)

	var bundleResponse interface{}

	resp, err := p.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var permission Permission
	err = permission.parseFromBundle(bundleResponse)
	return &permission, resp, err
}
