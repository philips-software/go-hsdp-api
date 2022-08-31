package iam

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// SMSTemplatesService represents the SMS template related services for IAM
type SMSTemplatesService struct {
	client *Client

	validate *validator.Validate
}

const (
	TypePhoneVerification      = "PHONE_VERIFICATION"
	TypeLoginOTP               = "LOGIN_OTP"
	TypePasswordRecovery       = "PASSWORD_RECOVERY"
	TypePasswordFailedAttempts = "PASSWORD_FAILED_ATTEMPTS"
)

type SMSTemplate struct {
	Schemas      []string          `json:"schemas" validate:"required"`
	ID           string            `json:"id,omitempty"`
	Organization OrganizationValue `json:"organization" validate:"required"`
	ExternalID   string            `json:"externalId,omitempty"`
	Type         string            `json:"type" validate:"required,oneof=PHONE_VERIFICATION MFA_OTP PASSWORD_RECOVERY PASSWORD_FAILED_ATTEMPTS"`
	Message      string            `json:"message" validate:"required"`
	Locale       string            `json:"locale,omitempty"`
	Meta         *Meta             `json:"meta,omitempty"`
}

// GetSMSTemplateOptions describes the criteria for looking up SMS templates
type GetSMSTemplateOptions struct {
	Filter             *string `url:"filter,omitempty"`
	Attributes         *string `url:"attributes,omitempty"`
	ExcludedAttributes *string `url:"excludedAttributes,omitempty"`
}

func SMSTemplateFilterOrgTypeLang(orgID, templateType, locale string) *GetSMSTemplateOptions {
	query := "organization.value eq \"" + orgID + "\" and type eq \"" + templateType + "\" and locale eq \"" + locale + "\""
	attributes := "id"
	return &GetSMSTemplateOptions{
		Filter:     &query,
		Attributes: &attributes,
	}
}

// CreateSMSTemplate creates a SMS template for IAM
func (o *SMSTemplatesService) CreateSMSTemplate(template SMSTemplate) (*SMSTemplate, *Response, error) {
	template.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSTemplate",
	}
	if err := o.validate.Struct(template); err != nil {
		return nil, nil, err
	}

	req, err := o.client.newRequest(IDM, "POST", "authorize/scim/v2/Configurations/SMSTemplate", &template, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var newTemplate SMSTemplate

	resp, err := o.client.do(req, &newTemplate)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, resp, fmt.Errorf("error creating sms template: %d", resp.StatusCode())
	}
	return &newTemplate, resp, err
}

// DeleteSMSTemplate deletes the SMS template
func (o *SMSTemplatesService) DeleteSMSTemplate(template SMSTemplate) (bool, *Response, error) {
	req, err := o.client.newRequest(IDM, "DELETE", "authorize/scim/v2/Configurations/SMSTemplate/"+template.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Method", "DELETE")

	var deleteResponse bytes.Buffer

	resp, err := o.client.do(req, &deleteResponse)
	if err != nil {
		return false, resp, err
	}
	return resp.StatusCode() == http.StatusAccepted, resp, nil
}

// UpdateSMSTemplate updates the SMS template
func (o *SMSTemplatesService) UpdateSMSTemplate(template SMSTemplate) (*SMSTemplate, *Response, error) {
	template.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSTemplate",
	}
	req, err := o.client.newRequest(IDM, "PUT", "authorize/scim/v2/Configurations/SMSTemplate/"+template.ID, &template, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", template.Meta.Version)

	var updatedTemplate SMSTemplate

	resp, err := o.client.do(req, &updatedTemplate)
	if err != nil {
		return nil, resp, err
	}
	return &updatedTemplate, resp, err

}

// GetSMSTemplateByID retrieves an SMS template by ID
func (o *SMSTemplatesService) GetSMSTemplateByID(id string) (*SMSTemplate, *Response, error) {
	var foundTemplate SMSTemplate

	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Configurations/SMSTemplate/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.do(req, &foundTemplate)
	if err != nil {
		return nil, resp, err
	}
	return &foundTemplate, resp, nil
}

// GetSMSTemplate retrieves an organization based on the GetSMSTemplateOptions parameters.
func (o *SMSTemplatesService) GetSMSTemplate(opt *GetSMSTemplateOptions, options ...OptionFunc) (*SMSTemplate, *Response, error) {
	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Configurations/SMSTemplate", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)

	var bundleResponse struct {
		Resources []struct {
			ID string `json:"id"`
		}
	}
	resp, err := o.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if len(bundleResponse.Resources) == 0 {
		return nil, resp, ErrNotFound
	}

	return o.GetSMSTemplateByID(bundleResponse.Resources[0].ID)
}
