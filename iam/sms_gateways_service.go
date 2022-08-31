package iam

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	smsServicesAPIVersion = "1"
)

// SMSGatewaysService represents the SMS related services for IAM
type SMSGatewaysService struct {
	client *Client

	validate *validator.Validate
}

type ProviderProperties struct {
	SID        string `json:"sid" validate:"required"`
	Endpoint   string `json:"endpoint" validate:"required"`
	FromNumber string `json:"fromNumber" validate:"required"`
}

type ProviderCredentials struct {
	Token string `json:"token" validate:"required"`
}

type OrganizationValue struct {
	Value string `json:"value" validate:"required"`
}

// SMSGateway represents a IAM SMS gateway resource
type SMSGateway struct {
	Schemas          []string            `json:"schemas" validate:"required"`
	ID               string              `json:"id,omitempty"`
	Organization     OrganizationValue   `json:"organization" validate:"required"`
	ExternalID       string              `json:"externalId,omitempty"`
	Provider         string              `json:"provider" validate:"required,oneof=twilio"`
	Properties       ProviderProperties  `json:"properties"`
	Credentials      ProviderCredentials `json:"credentials" validate:"required"`
	Active           bool                `json:"active"`
	ActivationExpiry int                 `json:"activationExpiry" validate:"required,min=1,max=43200"`
	Meta             *Meta               `json:"meta,omitempty"`
}

// GetSMSGatewayOptions describes the criteria for looking up SMS gateways
type GetSMSGatewayOptions struct {
	Filter             *string `url:"filter,omitempty"`
	Attributes         *string `url:"attributes,omitempty"`
	ExcludedAttributes *string `url:"excludedAttributes,omitempty"`
}

func SMSGatewayFilterOrgEq(orgID string) *GetSMSGatewayOptions {
	query := "id eq \"" + orgID + "\""
	attributes := "id"
	return &GetSMSGatewayOptions{
		Filter:     &query,
		Attributes: &attributes,
	}
}

// CreateSMSGateway creates a SMS gateway for IAM
func (o *SMSGatewaysService) CreateSMSGateway(gw SMSGateway) (*SMSGateway, *Response, error) {
	gw.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSGateway",
	}
	if err := o.validate.Struct(gw); err != nil {
		return nil, nil, err
	}

	req, err := o.client.newRequest(IDM, "POST", "authorize/scim/v2/Configurations/SMSGateway", &gw, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var newGW SMSGateway

	resp, err := o.client.do(req, &newGW)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, resp, fmt.Errorf("error creating sms gateway: %d", resp.StatusCode())
	}
	return &newGW, resp, err
}

// DeleteSMSGateway deletes the SMS gateway
func (o *SMSGatewaysService) DeleteSMSGateway(gw SMSGateway) (bool, *Response, error) {
	req, err := o.client.newRequest(IDM, "DELETE", "authorize/scim/v2/Configurations/SMSGateway/"+gw.ID, nil, nil)
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

// UpdateSMSGateway updates the SMS gateway
func (o *SMSGatewaysService) UpdateSMSGateway(gw SMSGateway) (*SMSGateway, *Response, error) {
	gw.Schemas = []string{
		"urn:ietf:params:scim:schemas:core:philips:hsdp:2.0:SMSGateway",
	}
	req, err := o.client.newRequest(IDM, "PUT", "authorize/scim/v2/Configurations/SMSGateway/"+gw.ID, &gw, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", gw.Meta.Version)

	var updatedGW SMSGateway

	resp, err := o.client.do(req, &updatedGW)
	if err != nil {
		return nil, resp, err
	}
	return &updatedGW, resp, err

}

// GetSMSGatewayByID retrieves an SMS gateway by ID
func (o *SMSGatewaysService) GetSMSGatewayByID(id string) (*SMSGateway, *Response, error) {
	var foundGW SMSGateway

	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Configurations/SMSGateway/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", smsServicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.do(req, &foundGW)
	if err != nil {
		return nil, resp, err
	}
	return &foundGW, resp, nil
}

// GetSMSGateway retrieves an SMS gateway based on the GetSMSGatewayOptions parameters.
func (o *SMSGatewaysService) GetSMSGateway(opt *GetSMSGatewayOptions, options ...OptionFunc) (*SMSGateway, *Response, error) {
	req, err := o.client.newRequest(IDM, "GET", "authorize/scim/v2/Configurations/SMSGateway", opt, options)
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

	return o.GetSMSGatewayByID(bundleResponse.Resources[0].ID)
}
