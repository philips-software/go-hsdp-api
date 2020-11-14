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

// ServiceOptions
type ServiceOptions struct {
}

func (c *ServicesService) GetRootCA(options ...OptionFunc) (*x509.Certificate, *Response, error) {
	return c.getCA("core/pki/api/root/ca/pem", options...)
}

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

func (c *ServicesService) GetRootCRL(options ...OptionFunc) (*pkix.CertificateList, *Response, error) {
	return c.getCRL("core/pki/api/root/crl/pem", options...)
}

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
