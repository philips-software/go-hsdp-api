package iam

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	emailTemplateAPIVersion = "1"
)

// EmailTemplatesService provides operations on IAM email template resources
type EmailTemplatesService struct {
	client *Client

	validate *validator.Validate
}

// EmailTemplate describes an email template
type EmailTemplate struct {
	// ID is the UUID generated for a stored email template
	ID string `json:"id,omitempty"`

	// Type is the type of the email template
	Type string `json:"type" validate:"required" enum:"ACCOUNT_ALREADY_VERIFIED|ACCOUNT_UNLOCKED|ACCOUNT_VERIFICATION|MFA_DISABLED|MFA_ENABLED|PASSWORD_CHANGED|PASSWORD_EXPIRY|PASSWORD_FAILED_ATTEMPTS|PASSWORD_RECOVERY"`

	// ManagingOrganization is the Unique UUID of the organization under which the email template needs to be created.
	ManagingOrganization string `json:"managingOrganization" validate:"required"`

	// From is the sender field
	From string `json:"from,omitempty"`

	// Format is the template format. Must be HTML at this time
	Format string `json:"format" validate:"required" enum:"HTML"`

	// Locale is the locale for the email template. The locale is case insensitive
	Locale string `json:"locale,omitempty"`

	// Subject is the email subject
	Subject string `json:"subject" validate:"required,min=1,max=256"`

	// Message should contain the base64 encoded body of the email
	Message string `json:"message" validate:"required"`

	// Link is a clickable link according to the template type
	Link string `json:"link,omitempty"`

	// Meta contains additional metadata
	Meta *Meta `json:"meta,omitempty"`
}

// CreateTemplate creates an EmailTemplate
// A user with EMAILTEMPLATE.WRITE permission can create templates under the organization.
func (e *EmailTemplatesService) CreateTemplate(template EmailTemplate) (*EmailTemplate, *Response, error) {
	if err := e.client.validate.Struct(template); err != nil {
		return nil, nil, err
	}
	req, err := e.client.newRequest(IDM, "POST", "authorize/identity/EmailTemplate", &template, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var createdTemplate EmailTemplate

	resp, err := e.client.do(req, &createdTemplate)
	if err != nil {
		return nil, resp, err
	}
	return &createdTemplate, resp, err
}

// DeleteTemplate deletes the given EmailTemplate
func (e *EmailTemplatesService) DeleteTemplate(template EmailTemplate) (bool, *Response, error) {
	req, err := e.client.newRequest(IDM, "DELETE", "authorize/identity/EmailTemplate/"+template.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var deleteResponse interface{}

	resp, err := e.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

type GetEmailTemplatesOptions struct {
	Type           *string `url:"type,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	Locale         *string `url:"locale,omitempty"`
}

// GetTemplates finds EmailTemplate based on search criteria
// Any user with EMAILTEMPLATE.WRITE or EMAILTEMPLATE.READ permission can retrieve the template information.
func (e *EmailTemplatesService) GetTemplates(opt *GetEmailTemplatesOptions, options ...OptionFunc) (*[]EmailTemplate, *Response, error) {
	req, err := e.client.newRequest(IDM, "GET", "authorize/identity/EmailTemplate", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var bundleResponse struct {
		Total int `json:"total"`
		Entry []struct {
			ID string `json:"id"`
		} `json:"entry"`
	}
	var templates []EmailTemplate

	resp, err := e.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrNotFound
	}
	for _, t := range bundleResponse.Entry {
		template, _, err := e.GetTemplateByID(t.ID)
		if err != nil {
			continue
		}
		templates = append(templates, *template)
	}
	return &templates, resp, nil
}

func (e *EmailTemplatesService) GetTemplateByID(ID string) (*EmailTemplate, *Response, error) {
	req, err := e.client.newRequest(IDM, "GET", "authorize/identity/EmailTemplate/"+ID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var template EmailTemplate

	resp, err := e.client.do(req, &template)
	if err != nil {
		return nil, resp, err
	}
	return &template, resp, err
}
