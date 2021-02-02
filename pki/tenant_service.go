package pki

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	AllowedDomains       []string `json:"allowed_domains,omitempty"`
	AllowedOtherSans     []string `json:"allowed_other_sans,omitempty"`
	AllowedSerialNumbers []string `json:"allowed_serial_numbers,omitempty"`
	AllowedURISans       []string `json:"allowed_uri_sans,omitempty"`
	ClientFlag           bool     `json:"client_flag" validate:"required"`
	Country              []string `json:"country"`
	EnforceHostnames     bool     `json:"enforce_hostnames" validate:"required"`
	KeyBits              int      `json:"key_bits"`
	KeyType              string   `json:"key_type"`
	Locality             []string `json:"locality,omitempty"`
	MaxTTL               string   `json:"max_ttl"`
	NotBeforeDuration    string   `json:"not_before_duration"`
	Organization         []string `json:"organization,omitempty"`
	OU                   []string `json:"ou,omitempty"`
	PostalCode           []string `json:"postal_code,omitempty"`
	Province             []string `json:"province,omitempty"`
	ServerFlag           bool     `json:"server_flag"`
	StreetAddress        []string `json:"street_address,omitempty"`
	TTL                  string   `json:"ttl"`
	UseCSRCommonName     bool     `json:"use_csr_common_name"`
	UseCSRSans           bool     `json:"use_csr_sans"`
}

type CertificateAuthority struct {
	TTL          string `json:"ttl"`
	CommonName   string `json:"common_name" validate:"required"`
	KeyType      string `json:"key_type"` // rsa|ec
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

type OnboardingResponse struct {
	APIEndpoint APIEndpoint `json:"api_endpoint"`
}

type APIEndpoint string

// LogicalPath returns the logical path component from the APIEndpoint
func (a APIEndpoint) LogicalPath() (string, error) {
	var logicalPath string
	u, err := url.Parse(string(a))
	if err != nil {
		return "", err
	}
	_, err = fmt.Sscanf(u.Path, "/core/pki/tenant/%s", &logicalPath)
	return logicalPath, err
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
	req, err := t.client.newTenantRequest(http.MethodPost, "core/pki/tenant", &tenant, options)
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
	resp, err := t.client.do(req, &onboardResponse)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	return &onboardResponse.OnboardingResponse, resp, nil
}

func (t *TenantService) Retrieve(logicalPath string, options ...OptionFunc) (*Tenant, *Response, error) {
	req, err := t.client.newTenantRequest(http.MethodGet, "core/pki/tenant/"+logicalPath, nil, options)
	if err != nil {
		return nil, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return nil, nil, err
	}
	var tenant Tenant
	resp, err := t.client.do(req, &tenant)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	return &tenant, resp, err
}

func (t *TenantService) Update(tenant Tenant, options ...OptionFunc) (bool, *Response, error) {
	if err := t.validate.Struct(tenant); err != nil {
		return false, nil, err
	}
	req, err := t.client.newTenantRequest(http.MethodPut, "core/pki/tenant/"+tenant.ServiceParameters.LogicalPath, &tenant, options)
	if err != nil {
		return false, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return false, nil, err
	}
	var errorResponse ErrorResponse
	resp, err := t.client.do(req, &errorResponse)
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
	req, err := t.client.newTenantRequest(http.MethodDelete, "core/pki/tenant/"+tenant.ServiceParameters.LogicalPath, &tenant, options)
	if err != nil {
		return false, nil, err
	}
	if err := t.setCFAuth(req); err != nil {
		return false, nil, err
	}
	var errorResponse ErrorResponse
	resp, err := t.client.do(req, &errorResponse)
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
