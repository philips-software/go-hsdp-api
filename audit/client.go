// Package audit provides support for HSDP Audit service
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"go.elastic.co/apm/module/apmhttp"

	signer "github.com/philips-software/go-hsdp-signer"

	"github.com/google/fhir/go/jsonformat"
)

const (
	libraryVersion = "0.29.0"
	userAgent      = "go-hsdp-api/audit/" + libraryVersion
	APIVersion     = "2"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region      string
	Environment string
	// ProductKey is provided as part of Auditing onboarding
	ProductKey string
	// Tenant value is used to support multi tenancy with a single ProductKey
	Tenant string
	// AuditBaseURL is provided as part of Auditing onboarding
	AuditBaseURL string
	// SharedKey is the IAM API signing key
	SharedKey string
	// SharedSecret is the IAM API signing secret
	SharedSecret string
	TimeZone     string
	DebugLog     string
}

// A Client manages communication with HSDP CDR API
type Client struct {
	config        *Config
	httpClient    *http.Client
	auditStoreURL *url.URL

	// User agent used when communicating with the HSDP IAM API.
	UserAgent string

	ma         *jsonformat.Marshaller
	um         *jsonformat.Unmarshaller
	httpSigner *signer.Signer

	debugFile *os.File
}

// NewClient returns a new HSDP Audit API client. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	return newClient(httpClient, config)
}

func newClient(httpClient *http.Client, config *Config) (*Client, error) {
	var err error

	if httpClient == nil {
		httpClient = apmhttp.WrapClient(http.DefaultClient)
	}

	c := &Client{httpClient: httpClient, config: config, UserAgent: userAgent}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		}
	}
	c.httpSigner, err = signer.New(c.config.SharedKey, c.config.SharedSecret)
	if err != nil {
		return nil, fmt.Errorf("signer.New: %w", err)
	}
	ma, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 marshaller: %w", err)
	}
	c.ma = ma
	um, err := jsonformat.NewUnmarshaller(config.TimeZone, jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 unmarshaller (timezone=[%s]): %w", config.TimeZone, err)
	}
	c.um = um
	c.setAuditBaseURL(c.config.AuditBaseURL)

	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

// setAuditBaseURL sets the FHIR store URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) setAuditBaseURL(urlStr string) error {
	if urlStr == "" {
		return ErrBaseURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.auditStoreURL, err = url.Parse(urlStr)
	return err
}

// newAuditRequest creates an new CDR Service API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newAuditRequest(method, path string, bodyBytes []byte, options []OptionFunc) (*http.Request, error) {
	u := *c.auditStoreURL
	// Set the encoded opaque data
	u.Path = c.auditStoreURL.Path + path

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

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("API-Version", APIVersion)
	req.Header.Set("Content-Type", "application/json")

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

// do executes a http request. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	if c.debugFile != nil {
		dumped, _ := httputil.DumpRequest(req, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Request start ---\n%s\n[go-hsdp-api] Request end ---\n", string(dumped))
		_, _ = c.debugFile.WriteString(out)
	}
	resp, err := c.httpClient.Do(req)
	if c.debugFile != nil && resp != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response start ---\n%s\n[go-hsdp-api] --- Response end ---\n", string(dumped))
		_, _ = c.debugFile.WriteString(out)
	}
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	doErr := CheckResponse(resp)

	if v != nil {
		defer resp.Body.Close() // Only close if we plan to read it
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return response, err
		}
	}

	return response, doErr
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	case 400:
		return ErrBadRequest
	}
	return ErrNonHttp20xResponse
}
