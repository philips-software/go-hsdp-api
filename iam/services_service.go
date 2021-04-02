package iam

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const servicesAPIVersion = "1"

// Service represents a IAM service resource
type Service struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name"`
	Description    string   `json:"description"` // RITM0021326
	ApplicationID  string   `json:"applicationId"`
	Validity       int      `json:"validity,omitempty"`
	ServiceID      string   `json:"serviceId,omitempty"`
	OrganizationID string   `json:"organizationId,omitempty"`
	ExpiresOn      string   `json:"expiresOn,omitempty"`
	PrivateKey     string   `json:"privateKey,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	DefaultScopes  []string `json:"defaultScopes,omitempty"`
}

// ServicesService provides operations on IAM Sessions resources
type ServicesService struct {
	client *Client
}

// GetServiceOptions describes search criteria for looking up services
type GetServiceOptions struct {
	ID             *string `url:"_id,omitempty"`
	Name           *string `url:"name,omitempty"`
	ApplicationID  *string `url:"applicationId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
	ServiceID      *string `url:"serviceId,omitempty"`
}

func fixHSDPPEM(pemString string) string {
	pre := strings.Replace(pemString,
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----\n", -1)
	return strings.Replace(pre,
		"-----END RSA PRIVATE KEY-----",
		"\n-----END RSA PRIVATE KEY-----", -1)
}

// GetToken returns a JWT which can be exchanged for an access token
func (s *Service) GetToken(accessTokenEndpoint string) (string, error) {
	// Decode private key
	block, _ := pem.Decode([]byte(fixHSDPPEM(s.PrivateKey)))
	if block == nil {
		return "", fmt.Errorf("failed to parse privateKey")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": accessTokenEndpoint,
		"iss": s.ServiceID,
		"sub": s.ServiceID,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	})
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// GetServiceByID looks up a service by ID
func (p *ServicesService) GetServiceByID(id string) (*Service, *Response, error) {
	return p.GetService(&GetServiceOptions{ID: &id}, nil)
}

// GetServiceByName looks up a service by name
func (p *ServicesService) GetServiceByName(name string) (*Service, *Response, error) {
	return p.GetService(&GetServiceOptions{Name: &name}, nil)
}

// GetServicesByApplicationID finds all services which belong to the applicationID
func (p *ServicesService) GetServicesByApplicationID(applicationID string) (*[]Service, *Response, error) {
	opt := &GetServiceOptions{
		ApplicationID: String(applicationID),
	}
	req, err := p.client.newRequest(IDM, "GET", "authorize/identity/Service", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)

	var responseStruct struct {
		Total int       `json:"total"`
		Entry []Service `json:"entry"`
	}

	resp, err := p.client.do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err
}

// CreateService creates a Service
func (p *ServicesService) CreateService(service Service) (*Service, *Response, error) {
	req, _ := p.client.newRequest(IDM, "POST", "authorize/identity/Service", &service, nil)
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var createdService Service

	resp, err := p.client.do(req, &createdService)
	if err != nil {
		return nil, resp, err
	}
	return &createdService, resp, err
}

// GetService looks up a services based on GetServiceOptions
func (p *ServicesService) GetService(opt *GetServiceOptions, options ...OptionFunc) (*Service, *Response, error) {
	services, resp, err := p.GetServices(opt, options...)
	if err != nil {
		return nil, resp, err
	}
	if len(*services) == 0 {
		return nil, resp, ErrEmptyResults
	}
	return &(*services)[0], resp, nil
}

// GetServices looks up services based on GetServiceOptions
func (p *ServicesService) GetServices(opt *GetServiceOptions, options ...OptionFunc) (*[]Service, *Response, error) {
	req, err := p.client.newRequest(IDM, "GET", "authorize/identity/Service", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int       `json:"total"`
		Entry []Service `json:"entry"`
	}

	resp, err := p.client.do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &bundleResponse.Entry, resp, err
}

// DeleteService deletes the given Service
func (p *ServicesService) DeleteService(service Service) (bool, *Response, error) {
	req, err := p.client.newRequest(IDM, "DELETE", "authorize/identity/Service/"+service.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var deleteResponse interface{}

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}

// AddScopes add scopes to the service
func (p *ServicesService) AddScopes(service Service, scopes []string, defaultScopes []string) (bool, *Response, error) {
	return p.updateScopes(service, "add", scopes, defaultScopes)
}

// RemoveScopes add scopes to the service
func (p *ServicesService) RemoveScopes(service Service, scopes []string, defaultScopes []string) (bool, *Response, error) {
	return p.updateScopes(service, "remove", scopes, defaultScopes)
}

func (p *ServicesService) updateScopes(service Service, action string, scopes []string, defaultScopes []string) (bool, *Response, error) {
	var requestBody = struct {
		Action        string   `json:"action"`
		Scopes        []string `json:"scopes,omitempty"`
		DefaultScopes []string `json:"defaultScopes,omitempty"`
	}{
		action,
		scopes,
		defaultScopes,
	}
	req, err := p.client.newRequest(IDM, "PUT", "authorize/identity/Service/"+service.ID+"/$scopes", requestBody, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var putResponse bytes.Buffer

	resp, err := p.client.do(req, &putResponse)
	if err != nil {
		return false, resp, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return false, resp, ErrOperationFailed
	}
	return true, resp, nil
}
