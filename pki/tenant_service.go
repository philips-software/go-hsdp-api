package pki

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type TenantService struct {
	orgID  string
	client *Client

	validate *validator.Validate
}

type Role struct {
	Name                 string   `json:"name" validate:"required"`
	AllowAnyName         bool     `json:"allow_any_name" validate:"required"`
	AllowIPSans          bool     `json:"allow_ip_sans" validate:"required"`
	AllowSubdomains      bool     `json:"allow_subdomains" validate:"required"`
	AllowedDomains       []string `json:"allowed_domains"`
	AllowedOtherSans     []string `json:"allowed_other_sans"`
	AllowedSerialNumbers []string `json:"allowed_serial_numbers"`
	AllowedURISans       []string `json:"allowed_uri_sans"`
	ClientFlag           bool     `json:"client_flag" validate:"required"`
	Country              []string `json:"country"`
	EnforceHostnames     bool     `json:"enforce_hostnames" validate:"required"`
	KeyBits              int      `json:"key_bits"`
	KeyType              string   `json:"key_type"`
	Locality             []string `json:"locality"`
	MaxTTL               string   `json:"max_ttl"`
	NotBeforeDuration    string   `json:"not_before_duration"`
	Organization         []string `json:"organization"`
	OU                   []string `json:"ou"`
	PostalCode           []string `json:"postal_code"`
	Province             []string `json:"province"`
	ServerFlag           bool     `json:"server_flag"`
	StreetAddress        []string `json:"street_address"`
	TTL                  string   `json:"ttl"`
	UseCSRCommonName     bool     `json:"use_csr_common_name"`
	UseCSRSans           bool     `json:"use_csr_sans"`
}

type CertificateAuthority struct {
	TTL          string `json:"ttl"`
	CommonName   string `json:"common_name" validate:"required"`
	KeyType      string `json:"key_type" validate:"required" enum:"rsa|ec"`
	KeyBits      int    `json:"key_bits"`
	OU           string `json:"ou"`
	Organization string `json:"organization"`
	Country      string `json:"country"`
	Locality     string `json:"locality"`
	Province     string `json:"province"`
}

type ServiceParameters struct {
	LogicalPath string               `json:"logical_path,omitempty"`
	IAMOrgs     []string             `json:"iam_orgs" validate:"min=1,max=10,required"`
	CA          CertificateAuthority `json:"ca" validate:"required"`
	Roles       []Role               `json:"roles" validate:"min=1,max=10,required"`
}

type Tenant struct {
	OrganizationName  string            `json:"organization_name" validate:"required"`
	SpaceName         string            `json:"space_name" validate:"required"`
	ServiceName       string            `json:"service_name" validate:"required"`
	PlanName          string            `json:"plan_name" validate:"required"`
	ServiceParameters ServiceParameters `json:"service_parameters" validate:"required"`
}

type ErrorResponse struct {
	Errors []string `json:"errors,omitempty"`
}

type OnboardingResponse struct {
	APIEndpoint string `json:"api_endpoint"`
}

func (t *TenantService) setCFAuth(req *http.Request) error {
	if t.client.consoleClient == nil {
		return ErrCFClientNotConfigured
	}
	token := t.client.consoleClient.Token()
	if token == "" {
		return ErrCFInvalidToken
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (t *TenantService) Onboard(tenant Tenant, options ...OptionFunc) (*OnboardingResponse, *Response, error) {
	if err := t.validate.Struct(tenant); err != nil {
		return nil, nil, err
	}
	req, err := t.client.NewTenantRequest(http.MethodPost, "core/pki/tenant", &tenant, options)
	if err != nil {
		return nil, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return nil, nil, err
	}
	var onboardResponse struct {
		ErrorResponse
		OnboardingResponse
	}
	resp, err := t.client.Do(req, &onboardResponse)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	return &onboardResponse.OnboardingResponse, resp, nil
}

func (t *TenantService) Retrieve(logicalPath string, options ...OptionFunc) (*Tenant, *Response, error) {
	req, err := t.client.NewTenantRequest(http.MethodGet, "core/pki/tenant/"+logicalPath, nil, options)
	if err != nil {
		return nil, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return nil, nil, err
	}
	var tenant Tenant
	resp, err := t.client.Do(req, &tenant)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	return &tenant, resp, err
}

func (t *TenantService) Update(tenant Tenant, options ...OptionFunc) (bool, *Response, error) {
	if err := t.validate.Struct(tenant); err != nil {
		return false, nil, err
	}
	req, err := t.client.NewTenantRequest(http.MethodPut, "core/pki/tenant/"+tenant.ServiceParameters.LogicalPath, &tenant, options)
	if err != nil {
		return false, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return false, nil, err
	}
	var errorResponse ErrorResponse
	resp, err := t.client.Do(req, &errorResponse)
	if err != nil && err != io.EOF {
		return false, nil, err
	}
	if resp == nil {
		return false, nil, ErrEmptyResult
	}
	if len(errorResponse.Errors) > 0 {
		err = fmt.Errorf("errors: %s", strings.Join(errorResponse.Errors, ","))
	} else {
		err = nil
	}
	return resp.StatusCode == http.StatusNoContent, resp, err
}

func (t *TenantService) Offboard(tenant Tenant, options ...OptionFunc) (bool, *Response, error) {
	req, err := t.client.NewTenantRequest(http.MethodDelete, "core/pki/tenant/"+tenant.ServiceParameters.LogicalPath, &tenant, options)
	if err != nil {
		return false, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return false, nil, err
	}
	var errorResponse ErrorResponse
	resp, err := t.client.Do(req, &errorResponse)
	if err != nil && err != io.EOF {
		return false, nil, err
	}
	if resp == nil {
		return false, nil, ErrEmptyResult
	}
	if len(errorResponse.Errors) > 0 {
		err = fmt.Errorf("errors: %s", strings.Join(errorResponse.Errors, ","))
	} else {
		err = nil
	}
	return resp.StatusCode == http.StatusNoContent, resp, nil
}
