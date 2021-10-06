package iam

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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

type CertificateOptionFunc func(cert *x509.Certificate) error

// FixPEM fixes the IAM generated PEM key strings so they are valid
// for decoding by Go and other parsers which expect newlines after labels
func FixPEM(pemString string) string {
	begin := "KEY-----"
	end := "-----END"
	pre := pemString
	if !strings.Contains(pre, begin+"\n") {
		pre = strings.Replace(pemString,
			begin,
			begin+"\n", -1)
	}
	if !strings.Contains(pre, "\n"+end) {
		return strings.Replace(pre,
			end,
			"\n"+end, -1)
	}
	return pre
}

// Valid checks if a service is usable
func (s *Service) Valid() bool {
	if s.ServiceID == "" || s.PrivateKey == "" {
		return false
	}
	return true
}

// GetToken returns a JWT which can be exchanged for an access token
func (s *Service) GetToken(accessTokenEndpoint string) (string, error) {
	// Decode private key
	block, _ := pem.Decode([]byte(FixPEM(s.PrivateKey)))
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

// UpdateServiceCertificateDER updates the associated certificate of the service using raw DER
func (p *ServicesService) UpdateServiceCertificateDER(service Service, derBytes []byte) (*Service, *Response, error) {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	var request = struct {
		Certificate string `json:"certificate"`
	}{
		Certificate: string(certPEM),
	}
	req, err := p.client.newRequest(IDM, "POST", "authorize/identity/Service/"+service.ID+"/$update-certificate", request, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	var updateResponse bytes.Buffer
	resp, err := p.client.do(req, &updateResponse)
	if err != nil {
		return nil, resp, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, resp, err
	}
	return p.GetServiceByID(service.ID)
}

// UpdateServiceCertificate updates the associated certificate of the service
func (p *ServicesService) UpdateServiceCertificate(service Service, privateKey *rsa.PrivateKey, options ...CertificateOptionFunc) (*Service, *Response, error) {
	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment
	notBefore := time.Now().Add(-24 * time.Hour)
	validFor := 365 * 24 * time.Hour
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template.Subject.CommonName = service.ServiceID
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	for _, o := range options {
		if err := o(&template); err != nil {
			return nil, nil, err
		}
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	if err != nil {
		return nil, nil, err
	}
	return p.UpdateServiceCertificateDER(service, derBytes)
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

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
