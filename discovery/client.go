// Package discovery provides support the HSDP Discovery services
package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/go-hsdp-api/internal"
)

const (
	userAgent  = "go-hsdp-api/discovery/" + internal.LibraryVersion
	APIVersion = "1"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a Client
type Config struct {
	Region      string
	Environment string
	BaseURL     string
	DebugLog    io.Writer
	Retry       int
}

// A Client manages communication with HSDP AI APIs
type Client struct {
	// HTTP Client used to communicate with IAM API
	*iam.Client
	config  *Config
	baseURL *url.URL

	// User agent used when communicating with the HSDP Notification API
	UserAgent string

	validate *validator.Validate
}

// Service describes a discovered service
type Service struct {
	ResourceType string   `json:"resourceType"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Tag          string   `json:"tag"`
	Actions      []string `json:"actions"`
	IsTrusted    bool     `json:"isTrusted"`
	URLS         []string `json:"urls"`
}

// NewClient returns a new Connect base Client
func NewClient(iamClient *iam.Client, config *Config) (*Client, error) {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}
	doAutoconf(config)
	c := &Client{Client: iamClient, config: config, UserAgent: userAgent, validate: validator.New()}

	if err := c.SetBaseURL(config.BaseURL); err != nil {
		return nil, err
	}

	return c, nil
}

func doAutoconf(config *Config) {
	if config.Region != "" && config.Environment != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region),
			autoconf.WithEnv(config.Environment))
		if err == nil {
			theService := c.Service("discovery")
			if theService.URL != "" && config.BaseURL == "" {
				config.BaseURL = theService.URL
			}
		}
	}
}

// Close releases allocated resources of clients
func (c *Client) Close() {
}

// GetBaseURL returns the base URL as configured
func (c *Client) GetBaseURL() string {
	if c.baseURL == nil {
		return ""
	}
	return c.baseURL.String()
}

// SetBaseURL sets the base URL for API requests
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

// GetEndpointURL returns the FHIR Store Endpoint URL as configured
func (c *Client) GetEndpointURL() string {
	return c.GetBaseURL()
}

func (c *Client) NewRequest(method, requestPath string, opt interface{}, options ...OptionFunc) (*http.Request, error) {
	u := *c.baseURL
	// Set the encoded opaque data
	u.Opaque = path.Join(c.baseURL.Path, requestPath)

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = io.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("API-Version", APIVersion)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// Response is a HSDP IAM API response. This wraps the standard http.Response
// returned from HSDP IAM and provides convenient access to things like errors
type Response struct {
	*http.Response
}

func (r *Response) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode
	}
	return 0
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
	resp, err := c.HttpClient().Do(req)
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
		defer func() {
			_ = resp.Body.Close()
		}() // Only close if we plan to read it
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

func (c *Client) GetServices() (*[]Service, *Response, error) {
	requiredScope := "?.?.dsc.service.readAny"
	if !c.HasScopes(requiredScope) {
		return nil, nil, fmt.Errorf("missing scope '%s'", requiredScope)
	}

	req, err := c.NewRequest(http.MethodGet, "Service", nil)
	if err != nil {
		return nil, nil, err
	}
	var bundleResponse internal.Bundle

	resp, err := c.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if err := internal.CheckResponse(resp.Response); err != nil {
		return nil, resp, err
	}
	var services []Service
	for _, s := range bundleResponse.Entry {
		var service Service
		if err := json.Unmarshal(s.Resource, &service); err == nil {
			services = append(services, service)
		}
	}
	return &services, resp, nil
}
