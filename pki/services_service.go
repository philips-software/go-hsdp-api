package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ServicesService struct {
	orgID  string
	client *Client

	validate *validator.Validate
}

type CertificateRequest struct {
	CommonName        string `json:"common_name" validate:"required,max=253"`
	AltName           string `json:"alt_name,omitempty"`
	IPSANS            string `json:"ip_sans,omitempty"`
	URISANS           string `json:"uri_sans,omitempty"`
	OtherSANS         string `json:"other_sans,omitempty"`
	TTL               string `json:"ttl,omitempty"`
	Format            string `json:"format,omitempty"`
	PrivateKeyFormat  string `json:"private_key_format,omitempty"`
	ExcludeCNFromSANS *bool  `json:"exclude_cn_from_sans,omitempty"`
}

type IssueData struct {
	CaChain        []string `json:"ca_chain"`
	Certificate    string   `json:"certificate"`
	Expiration     int      `json:"expiration"`
	IssuingCa      string   `json:"issuing_ca"`
	PrivateKey     string   `json:"private_key"`
	PrivateKeyType string   `json:"private_key_type"`
	SerialNumber   string   `json:"serial_number"`
}

type IssueResponse struct {
	RequestID     string    `json:"request_id"`
	LeaseID       string    `json:"lease_id"`
	Renewable     bool      `json:"renewable"`
	LeaseDuration int       `json:"lease_duration"`
	Data          IssueData `json:"data"`
	WrapInfo      *string   `json:"wrap_info"`
	Warnings      *string   `json:"warnings"`
	Auth          *string   `json:"auth"`
}

// ServiceOptions
type ServiceOptions struct {
}

// GetRootCA
func (c *ServicesService) GetRootCA(options ...OptionFunc) (*x509.Certificate, *Response, error) {
	return c.getCA("core/pki/api/root/ca/pem", options...)
}

// GetPolicyCA
func (c *ServicesService) GetPolicyCA(options ...OptionFunc) (*x509.Certificate, *Response, error) {
	return c.getCA("core/pki/api/policy/ca/pem", options...)
}

func (c *ServicesService) getCA(path string, options ...OptionFunc) (*x509.Certificate, *Response, error) {
	req, err := c.client.NewServiceRequest(http.MethodGet, path, nil, options)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.client.Do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	defer resp.Body.Close()
	pemData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, resp, ErrCertificateExpected
	}
	pub, err := x509.ParseCertificate(block.Bytes)
	return pub, resp, err
}

// GetRootCRL
func (c *ServicesService) GetRootCRL(options ...OptionFunc) (*pkix.CertificateList, *Response, error) {
	return c.getCRL("core/pki/api/root/crl/pem", options...)
}

// GetPolicyCRL
func (c *ServicesService) GetPolicyCRL(options ...OptionFunc) (*pkix.CertificateList, *Response, error) {
	return c.getCRL("core/pki/api/policy/crl/pem", options...)
}

func (c *ServicesService) getCRL(path string, options ...OptionFunc) (*pkix.CertificateList, *Response, error) {
	req, err := c.client.NewServiceRequest(http.MethodGet, path, nil, options)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.client.Do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	defer resp.Body.Close()
	pemData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "X509 CRL" {
		return nil, resp, ErrCRLExpected
	}
	pub, err := x509.ParseCRL(block.Bytes)
	return pub, resp, err
}

func (d *IssueData) GetCertificate() (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(d.Certificate))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, ErrCertificateExpected
	}
	return x509.ParseCertificate(block.Bytes)
}

func (d *IssueData) GetPrivateKey() (interface{}, error) {
	block, _ := pem.Decode([]byte(d.PrivateKey))
	if block == nil {
		return nil, ErrInvalidPrivateKey
	}
	switch d.PrivateKeyType {
	case "rsa":
		if block.Type != "RSA PRIVATE KEY" {
			return nil, ErrInvalidPrivateKey
		}
		private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return private, nil
	case "ec":
		if block.Type != "EC PRIVATE KEY" {
			return nil, ErrInvalidPrivateKey
		}
		private, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return private, nil
	}
	return nil, ErrInvalidPrivateKey
}

// IssueCertificate
func (c *ServicesService) IssueCertificate(logicalPath, roleName string, request CertificateRequest, options ...OptionFunc) (*IssueResponse, *Response, error) {
	req, err := c.client.NewServiceRequest(http.MethodPost, "core/pki/api/"+logicalPath+"/issue/"+roleName, &request, options)
	if err != nil {
		return nil, nil, err
	}
	var responseStruct struct {
		IssueResponse
		ErrorResponse
	}
	resp, err := c.client.Do(req, &responseStruct)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	return &responseStruct.IssueResponse, resp, nil
}

// GetCertificateBySerial
func (c *ServicesService) GetCertificateBySerial(logicalPath, serial string, options ...OptionFunc) (*IssueResponse, *Response, error) {
	req, err := c.client.NewServiceRequest(http.MethodGet, "core/pki/api/"+logicalPath+"/cert/"+serial, nil, options)
	var responseStruct struct {
		IssueResponse
		ErrorResponse
	}
	resp, err := c.client.Do(req, &responseStruct)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, ErrEmptyResult
	}
	return &responseStruct.IssueResponse, resp, nil
}
