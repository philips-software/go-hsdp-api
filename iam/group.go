package iam

// Group represents an IAM group resource
type Group struct {
	ID                   string `json:"id,omitempty" validate:""`
	Name                 string `json:"name,omitempty" validate:"required"`
	Description          string `json:"description,omitempty" validate:""`
	ManagingOrganization string `json:"managingOrganization,omitempty" validate:"required"`
}
