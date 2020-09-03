// Package cce provides support for CCE
package cce

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

	"github.com/google/go-querystring/query"
	"github.com/philips-software/go-hsdp-api/fhir"
	"github.com/philips-software/go-hsdp-api/iam"
)

const (
	libraryVersion = "0.23.0"
	userAgent      = "go-hsdp-api/cce/" + libraryVersion
	DiscoveryPath  = ".well-known/cce-configuration"
	iamTokenPath   = "/authorize/oauth2/token"
	APIVersion     = "1"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	client     *http.Client
	ServiceID  string
	PrivateKey string
	BaseURL    string
	Debug      bool
	DebugLog   string
}

type DiscoveryEndpoints struct {
	TokenEndpoint         string `json:"token_endpoint"`
	IntrospectionEndpoint string `json:"introspection_endpoint"`
	DiscoveryEndpoint     string `json:"discovery_endpoint"`
}

// A Client manages communication with HSDP IAM API
type Client struct {
	// HTTP client used to communicate with the API.
	iamClient *iam.Client
	client    *http.Client

	config *Config

	baseURL *url.URL

	Endpoints DiscoveryEndpoints

	// User agent used when communicating with the API server.
	UserAgent string

	debugFile *os.File
	/*
		Policy *PolicyService
		Access *AccessService
	*/
}

// NewClient returns a new CCE API client. A configured IAM
// client must be provided
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	return newClient(httpClient, config)
}

func newClient(httpClient *http.Client, config *Config) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient, config: config, UserAgent: userAgent}
	if err := c.SetBaseURL(c.config.BaseURL); err != nil {
		return nil, err
	}
	// Debug
	if c.config.DebugLog != "" {
		var err error
		debugFile, err := os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		c.debugFile = debugFile
	}

	// Endpoints
	if err := c.Bootstrap(); err != nil {
		return nil, err
	}

	// IAM token
	iamClient, err := iam.NewClient(config.client, &iam.Config{
		IAMURL:   strings.TrimSuffix(c.Endpoints.TokenEndpoint, iamTokenPath),
		Debug:    config.Debug,
		DebugLog: config.DebugLog,
	})
	if err != nil {
		return nil, err
	}
	c.iamClient = iamClient

	err = iamClient.ServiceLogin(iam.Service{
		ServiceID:  config.ServiceID,
		PrivateKey: config.PrivateKey,
	})
	if err != nil {
		return nil, err
	}

	/*
		c.Policy = &PolicyService{client: c, validate: validator.New()}
		_ = c.Policy.validate.RegisterValidation("policyActions", validateActions)

		c.Access = &AccessService{client: c}
	*/

	return c, nil
}

// Bootstrap does endpoint discovery of CCE
func (c *Client) Bootstrap() error {
	// Endpoints
	u := *c.baseURL
	// Set the encoded opaque data
	u.Opaque = c.baseURL.Path + DiscoveryPath
	req := &http.Request{
		Method:     "GET",
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	req.Header.Set("Api-Version", APIVersion)
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.Do(req, &c.Endpoints)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrDiscoveryFailed
	}
	return nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

// SetBaseURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseURL(urlStr string) error {
	if urlStr == "" {
		return ErrBaseURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseURL, err = url.Parse(urlStr)
	return err
}

// NewRequest creates an new CCE API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, endpoint string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req := &http.Request{
		Method:     method,
		URL:        u,
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

	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.iamClient.Token())

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response is a CCE API response. This wraps the standard http.Response
// returned from CCE and provides convenient access to things like errors
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Do executes a http request. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	if c.config.Debug {
		dumped, _ := httputil.DumpRequest(req, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Request start ---\n%s\n[go-hsdp-api] Request end ---\n", string(dumped))
		if c.debugFile != nil {
			_, _ = c.debugFile.WriteString(out)
		} else {
			fmt.Println(out)
		}
	}
	resp, err := c.client.Do(req)
	if c.config.Debug && resp != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response start ---\n%s\n[go-hsdp-api] --- Response end ---\n", string(dumped))
		if c.debugFile != nil {
			_, _ = c.debugFile.WriteString(out)
		} else {
			fmt.Println(out)
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = fhir.CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}
