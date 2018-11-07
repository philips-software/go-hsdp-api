package iam

// Application represents an IAM Application entity
type Application struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name" validate:"required"`
	Description       string `json:"description"`
	PropositionID     string `json:"propositionId" validate:"required"`
	GlobalReferenceID string `json:"globalReferenceId" validate:"required"`
}
