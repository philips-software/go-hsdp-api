// Package dicom provides support for HSDP DICOM services
package dicom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/philips-software/go-hsdp-api/internal"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/google/fhir/go/jsonformat"

	"github.com/philips-software/go-hsdp-api/iam"
)

const (
	userAgent  = "go-hsdp-api/dicom/" + internal.LibraryVersion
	APIVersion = "1"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region      string
	Environment string
	RootOrgID   string
	DICOMURL    string
	DICOMStore  string
	Type        string
	TimeZone    string
	DebugLog    string
}

// A Client manages communication with HSDP DICOM API
type Client struct {
	// HTTP client used to communicate with IAM API
	iamClient *iam.Client

	config *Config

	dicomStoreURL *url.URL

	// User agent used when communicating with the HSDP DICOM API.
	UserAgent string

	debugFile *os.File

	Config *ConfigService
}

// NewClient returns a new HSDP CDR API client. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(iamClient *iam.Client, config *Config) (*Client, error) {
	return newClient(iamClient, config)
}

func newClient(iamClient *iam.Client, config *Config) (*Client, error) {
	c := &Client{iamClient: iamClient, config: config, UserAgent: userAgent}
	dicomStore := config.DICOMStore
	if dicomStore == "" {
		dicomStore = config.DICOMURL + "/store/dicom/"
	}
	if err := c.SetDICOMStoreURL(dicomStore); err != nil {
		return nil, err
	}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		}
	}
	ma, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 marshaller: %w", err)
	}
	um, err := jsonformat.NewUnmarshaller(config.TimeZone, jsonformat.STU3)
	if err != nil {
		return nil, fmt.Errorf("cdr.NewClient create FHIR STU3 unmarshaller (timezone=[%s]): %w", config.TimeZone, err)
	}

	c.Config = &ConfigService{client: c, ma: ma, um: um, profile: "production"}

	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

// GetDICOMStoreURL returns the base FHIR Store base URL as configured
func (c *Client) GetDICOMStoreURL() string {
	if c.dicomStoreURL == nil {
		return ""
	}
	return c.dicomStoreURL.String()
}

// GetEndpointURL returns the FHIR Store Endpoint URL as configured
func (c *Client) GetEndpointURL() string {
	return c.GetDICOMStoreURL() + c.config.RootOrgID
}

// SetDICOMStoreURL sets the FHIR store URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetDICOMStoreURL(urlStr string) error {
	if urlStr == "" {
		return ErrDICOMURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.dicomStoreURL, err = url.Parse(urlStr)
	return err
}

// SetEndpointURL sets the FHIR endpoint URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetEndpointURL(urlStr string) error {
	if urlStr == "" {
		return ErrDICOMURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.dicomStoreURL, err = url.Parse(urlStr)
	if err != nil {
		return err
	}
	parts := strings.Split(c.dicomStoreURL.Path, "/")
	if len(parts) == 0 {
		return ErrDICOMURLCannotBeEmpty
	}
	c.config.RootOrgID = parts[len(parts)-1]
	newParts := parts[:len(parts)-1]
	c.dicomStoreURL.Path = strings.Join(newParts, "/")
	return nil
}

// newDICOMRequest creates an new DICOM Service API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newDICOMRequest(method, path string, bodyBytes []byte, options ...OptionFunc) (*http.Request, error) {
	u := *c.dicomStoreURL
	// Set the encoded opaque data
	u.Opaque = c.dicomStoreURL.Path + path

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
	req.Header.Set("Authorization", "Bearer "+c.iamClient.Token())
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

// do executes a http request. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	if c.debugFile != nil {
		dumped, _ := httputil.DumpRequest(req, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Request start ---\n%s\n[go-hsdp-api] Request end ---\n", string(dumped))
		_, _ = c.debugFile.WriteString(out)
	}
	resp, err := c.iamClient.HttpClient().Do(req)
	if c.debugFile != nil && resp != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response start ---\n%s\n[go-hsdp-api] --- Response end ---\n", string(dumped))
		_, _ = c.debugFile.WriteString(out)
	}
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	err = checkResponse(resp)
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

// checkResponse checks the API response for errors, and returns them if present.
func checkResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}
	return ErrNonHttp20xResponse
}
