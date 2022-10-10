// Package audit provides support for interacting with the HSDP Audit service
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/fhir/go/fhirversion"
	"github.com/philips-software/go-hsdp-api/internal"

	signer "github.com/philips-software/go-hsdp-signer"

	"github.com/google/fhir/go/jsonformat"
)

const (
	userAgent  = "go-hsdp-api/audit/" + internal.LibraryVersion
	APIVersion = "2"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region      string
	Environment string
	// AuditBaseURL is provided as part of Auditing onboarding
	AuditBaseURL string
	// SharedKey is the IAM API signing key
	SharedKey string
	// SharedSecret is the IAM API signing secret
	SharedSecret string
	TimeZone     string
	DebugLog     string
}

// Client holds state of a HSDP Audit client
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
		c := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
		httpClient = c
	}

	c := &Client{httpClient: httpClient, config: config, UserAgent: userAgent}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			httpClient.Transport = internal.NewLoggingRoundTripper(httpClient.Transport, c.debugFile)
		}
	}
	c.httpSigner, err = signer.New(c.config.SharedKey, c.config.SharedSecret)
	if err != nil {
		return nil, fmt.Errorf("signer.New: %w", err)
	}
	ma, err := jsonformat.NewMarshaller(false, "", "", fhirversion.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 marshaller: %w", err)
	}
	c.ma = ma
	um, err := jsonformat.NewUnmarshaller(config.TimeZone, fhirversion.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 unmarshaller (timezone=[%s]): %w", config.TimeZone, err)
	}
	c.um = um
	_ = c.setAuditBaseURL(c.config.AuditBaseURL)

	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

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
		req.Body = io.NopCloser(bodyReader)
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

// Response is a HSDP Audit API response. This wraps the standard http.Response
// returned from HSDP Audit and provides convenient access to things like errors
type Response struct {
	*http.Response
}

func (r *Response) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode
	}
	return 0
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	doErr := internal.CheckResponse(resp)

	if v != nil {
		defer func() {
			_ = resp.Body.Close()
		}() // Only close if we plan to read it
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
