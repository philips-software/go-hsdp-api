// Package cdr provides support for HSDP CDR services
package cdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/philips-software/go-hsdp-api/internal"

	"github.com/google/fhir/go/jsonformat"

	"github.com/philips-software/go-hsdp-api/iam"
)

const (
	userAgent  = "go-hsdp-api/cdr/" + internal.LibraryVersion
	APIVersion = "1"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region      string
	Environment string
	RootOrgID   string
	// CDRURL is the URL of the CDR instance, including the /store/fhir or /store/personal suffix path
	CDRURL    string
	FHIRStore string
	Type      string
	TimeZone  string
	DebugLog  string
}

// A Client manages communication with HSDP CDR API
type Client struct {
	// HTTP client used to communicate with IAM API
	iamClient *iam.Client

	config *Config

	fhirStoreURL *url.URL

	// User agent used when communicating with the HSDP CDR API
	UserAgent string

	TenantSTU3     *TenantSTU3Service
	OperationsSTU3 *OperationsSTU3Service

	TenantR4     *TenantR4Service
	OperationsR4 *OperationsR4Service
}

// NewClient returns a new HSDP CDR API client. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(iamClient *iam.Client, config *Config) (*Client, error) {
	return newClient(iamClient, config)
}

func newClient(iamClient *iam.Client, config *Config) (*Client, error) {
	c := &Client{iamClient: iamClient, config: config, UserAgent: userAgent}
	fhirStore := config.FHIRStore
	if fhirStore == "" {
		fhirStore = config.CDRURL
	}
	if err := c.SetFHIRStoreURL(fhirStore); err != nil {
		return nil, err
	}
	maSTU3, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 marshaller: %w", err)
	}
	umSTU3, err := jsonformat.NewUnmarshaller(config.TimeZone, jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 unmarshaller (timezone=[%s]): %w", config.TimeZone, err)
	}
	maR4, err := jsonformat.NewMarshaller(false, "", "", jsonformat.R4)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 marshaller: %w", err)
	}
	umR4, err := jsonformat.NewUnmarshaller(config.TimeZone, jsonformat.R4)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 unmarshaller (timezone=[%s]): %w", config.TimeZone, err)
	}

	c.TenantSTU3 = &TenantSTU3Service{timeZone: config.TimeZone, client: c, ma: maSTU3, um: umSTU3}
	c.OperationsSTU3 = &OperationsSTU3Service{timeZone: config.TimeZone, client: c, ma: maSTU3, um: umSTU3}
	c.TenantR4 = &TenantR4Service{timeZone: config.TimeZone, client: c, ma: maR4, um: umR4}
	c.OperationsR4 = &OperationsR4Service{timeZone: config.TimeZone, client: c, ma: maR4, um: umR4}

	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
}

// GetFHIRStoreURL returns the base FHIR Store base URL as configured
func (c *Client) GetFHIRStoreURL() string {
	if c.fhirStoreURL == nil {
		return ""
	}
	return c.fhirStoreURL.String()
}

// GetEndpointURL returns the FHIR Store Endpoint URL as configured
func (c *Client) GetEndpointURL() string {
	return c.GetFHIRStoreURL() + c.config.RootOrgID
}

// SetFHIRStoreURL sets the FHIR store URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetFHIRStoreURL(urlStr string) error {
	if urlStr == "" {
		return ErrCDRURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.fhirStoreURL, err = url.Parse(urlStr)
	return err
}

// SetEndpointURL sets the FHIR endpoint URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetEndpointURL(urlStr string) error {
	if urlStr == "" {
		return ErrCDRURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.fhirStoreURL, err = url.Parse(urlStr)
	if err != nil {
		return err
	}
	parts := strings.Split(c.fhirStoreURL.Path, "/")
	if len(parts) == 0 {
		return ErrCDRURLCannotBeEmpty
	}
	c.config.RootOrgID = parts[len(parts)-1]
	newParts := parts[:len(parts)-1]
	c.fhirStoreURL.Path = strings.Join(newParts, "/")
	return nil
}

// newCDRRequest creates an new CDR Service API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newCDRRequest(method, path string, bodyBytes []byte, options []OptionFunc) (*http.Request, error) {
	u := *c.fhirStoreURL
	// Set the encoded opaque data
	u.Opaque = c.fhirStoreURL.Path + c.config.RootOrgID + "/" + path

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}

		if err := fn(req); err != nil {
			return nil, err
		}
	}

	if method == "POST" || method == "PUT" || method == "PATCH" {
		bodyReader := bytes.NewReader(bodyBytes)
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
	}
	token, err := c.iamClient.Token()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("API-Version", APIVersion)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response is a HSDP IAM API response. This wraps the standard http.Response
// returned from HSDP IAM and provides convenient access to things like errors
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// TokenRefresh forces a refresh of the IAM access token
func (c *Client) TokenRefresh() error {
	if c.iamClient == nil {
		return fmt.Errorf("invalid IAM client, cannot refresh token")
	}
	return c.iamClient.TokenRefresh()
}

// do executes a http request. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	if req.Header.Get("Accept") == "" {
		return nil, ErrMissingAcceptHeader
	}

	resp, err := c.iamClient.HttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	err = internal.CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		defer resp.Body.Close() // Only close if we plan to read it
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}
