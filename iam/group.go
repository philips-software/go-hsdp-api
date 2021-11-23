package iam

// Group represents an IAM group resource
type Group struct {
	ID                   string `json:"id,omitempty" validate:""`
	Name                 string `json:"name,omitempty" validate:"required"`
	Description          string `json:"description,omitempty" validate:""`
	ManagingOrganization string `json:"managingOrganization,omitempty" validate:"required"`
}

// GroupResource is the resource response of a Group search operation
type GroupResource struct {
	ID               string `json:"_id"`
	ResourceType     string `json:"resourceType"`
	GroupName        string `json:"groupName"`
	OrgID            string `json:"orgId"`
	GroupDescription string `json:"groupDescription"`
}
