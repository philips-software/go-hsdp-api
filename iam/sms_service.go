package iam

const smsServicesAPIVersion = "1"

// SMSGateway represents a IAM SMS gateway resource
type SMSGateway struct {
	Schemas      []string `json:"schemas" validate:"required"`
	ID           string   `json:"id,omitempty"`
	Organization struct {
		Value string `json:"value" validate:"required"`
	} `json:"organization"`
	ExternalID string `json:"externalId,omitempty"`
	Provider   string `json:"provider" validate:"required"`
	Properties struct {
		SID        string `json:"sid" validate:"required"`
		Endpoint   string `json:"endpoint" validate:"required"`
		FromNumber string `json:"fromNumber" validate:"required"`
	} `json:"properties"`
	Credentials struct {
		Token string `json:"token" validate:"required"`
	} `json:"credentials" validate:"required"`
	Active           bool `json:"active" validate:"required"`
	ActivationExpiry int  `json:"activationExpiry" validate:"required,min=1,max=43200"`
}
