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

type EmailTemplate struct {
	ID                   string `json:"id,omitempty"`
	Type                 string `json:"type" validate:"required" enum:"ACCOUNT_ALREADY_VERIFIED|ACCOUNT_UNLOCKED|ACCOUNT_VERIFICATION|MFA_DISABLED|MFA_ENABLED|PASSWORD_CHANGED|PASSWORD_EXPIRY|PASSWORD_FAILED_ATTEMPTS|PASSWORD_RECOVERY"`
	ManagingOrganization string `json:"managingOrganization" validate:"required"`
	Format               string `json:"format" validate:"required" enum:"HTML"`
	Locale               string `json:"locale"  validate:"required"`
	Subject              string `json:"subject"`
	Message              string `json:"message" validate:"required"`
	Link                 string `json:"link"`
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

	resp, err := e.client.Do(req, &createdTemplate)
	if err != nil {
		return nil, resp, err
	}
	return &createdTemplate, resp, err
}

// DeleteTemplate deletes the given EmailTemplate
func (g *GroupsService) DeleteTemplate(template EmailTemplate) (bool, *Response, error) {
	req, err := g.client.newRequest(IDM, "DELETE", "authorize/identity/EmailTemplate/"+template.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var deleteResponse interface{}

	resp, err := g.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

type GetEmailTemplatesOptions struct {
	ID             *string `url:"_id,omitempty"`
	Name           *string `url:"name,omitempty"`
	ApplicationID  *string `url:"applicationId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	ServiceID      *string `url:"serviceId,omitempty"`
}

// GetTemplate finds EmailTemplate based on search criteria
// Any user with EMAILTEMPLATE.WRITE or EMAILTEMPLATE.READ permission can retrieve the template information.
func (e *EmailTemplatesService) GetTemplate(opt *GetEmailTemplatesOptions, options ...OptionFunc) (*EmailTemplate, *Response, error) {
	req, err := e.client.newRequest(IDM, "GET", "authorize/identity/EmailTemplate", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var bundleResponse struct {
		Total int `json:"total"`
		Entry []struct {
			Resource struct {
				ID string `json:"_id"`
			} `json:"resource"`
		} `json:"entry"`
	}

	resp, err := e.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return nil, resp, ErrNotFound
	}
	return e.GetTemplateByID(bundleResponse.Entry[0].Resource.ID)
}

func (e *EmailTemplatesService) GetTemplateByID(ID string) (*EmailTemplate, *Response, error) {
	req, err := e.client.newRequest(IDM, "GET", "authorize/identity/EmailTemplate/"+ID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", emailTemplateAPIVersion)

	var template EmailTemplate

	resp, err := e.client.Do(req, &template)
	if err != nil {
		return nil, resp, err
	}
	return &template, resp, err
}
