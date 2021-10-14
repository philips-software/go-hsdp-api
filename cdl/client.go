// Package cdl provides support for interacting with HSDP Clinical Data Lake services
package cdl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/go-hsdp-api/internal"
)

const (
	userAgent  = "go-hsdp-api/cdl/" + internal.LibraryVersion
	APIVersion = "3"
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region         string
	Environment    string
	OrganizationID string `validate:"required"`
	CDLURL         string
	CDLStore       string
	DebugLog       string
	Retry          int
}

// A Client manages communication with HSDP CDL API
type Client struct {
	// HTTP client used to communicate with IAM API
	iamClient *iam.Client

	config *Config

	cdlURL *url.URL

	cdlStoreURL *url.URL

	// User agent used when communicating with the HSDP Notification API
	UserAgent string

	debugFile *os.File
	validate  *validator.Validate

	Study              *StudyService
	DataTypeDefinition *DatatypeDefinitionService
	LabelDefinition    *LabelDefinitionService
	ExportRoute        *ExportRouteService
}

// NewClient returns a new HSDP CDL API client. A configured IAM client
// must be provided as the underlying API requires an IAM token
func NewClient(iamClient *iam.Client, config *Config) (*Client, error) {
	return newClient(iamClient, config)
}

func newClient(iamClient *iam.Client, config *Config) (*Client, error) {
	doAutoconf(config)
	c := &Client{iamClient: iamClient, config: config, UserAgent: userAgent, validate: validator.New()}
	cdlStore := config.CDLStore
	if cdlStore == "" {
		cdlStore = config.CDLURL + "/store/cdl/" + c.config.OrganizationID
	}
	if err := c.SetCDLStoreURL(cdlStore); err != nil {
		return nil, err
	}

	c.Study = &StudyService{client: c, validate: validator.New(), config: config}
	c.DataTypeDefinition = &DatatypeDefinitionService{client: c, validate: validator.New(), config: config}
	c.LabelDefinition = &LabelDefinitionService{client: c, validate: validator.New(), config: config}
	_ = c.LabelDefinition.validate.RegisterValidation("labelDefValidationHandler", labelDefValidationHandler)

	c.ExportRoute = &ExportRouteService{client: c, validate: validator.New(), config: config}
	return c, nil
}

func labelDefValidationHandler(f1 validator.FieldLevel) bool {
	switch f1.FieldName() {
	case "LabelDefName":
		return len(f1.Field().String()) > 0
	case "Description":
		return len(f1.Field().String()) != 0
	case "LabelScope":
		labelScope := f1.Field().Interface().(LabelScope)
		return len(labelScope.Type) != 0
	case "Label":
		return len(f1.Field().String()) != 0
	case "Type":
		return len(f1.Field().String()) != 0
	case "Labels":
		labelsArray := f1.Field().Interface().([]LabelsArrayElem)
		for _, i := range labelsArray {
			if len(i.Label) == 0 {
				return false
			}
		}
		return len(labelsArray) != 0
	default:
		return false
	}
}

func doAutoconf(config *Config) {
	if config.Region != "" && config.Environment != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region),
			autoconf.WithEnv(config.Environment))
		if err == nil {
			cdlService := c.Service("cdl")
			if cdlService.URL != "" && config.CDLURL == "" {
				config.CDLURL = cdlService.URL
			}
		}
	}
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

// GetCDLURL returns the base CDL Store base URL as configured
func (c *Client) GetCDLURL() string {
	if c.cdlURL == nil {
		return ""
	}
	return c.cdlURL.String()
}

// SetCDLURL sets the Notification URL for API requests
func (c *Client) SetCDLURL(urlStr string) error {
	if urlStr == "" {
		return ErrCDLURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.cdlURL, err = url.Parse(urlStr)
	return err
}

// GetEndpointURL returns the CDL Store URL including the tenant ID
func (c *Client) GetEndpointURL() string {
	return c.GetCDLStoreURL()
}

// SetEndpointURL sets the CDL endpoint URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetEndpointURL(urlStr string) error {
	if urlStr == "" {
		return ErrCDLURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.cdlStoreURL, err = url.Parse(urlStr)
	if err != nil {
		return err
	}
	parts := strings.Split(c.cdlStoreURL.Path, "/")
	if len(parts) == 0 {
		return ErrCDLURLCannotBeEmpty
	}
	c.config.OrganizationID = parts[len(parts)-1]
	newParts := parts[:len(parts)-1]
	c.cdlStoreURL.Path = strings.Join(newParts, "/") + "/"
	return nil
}

// GetCDLStoreURL returns the base CDL Store base URL as configured
func (c *Client) GetCDLStoreURL() string {
	if c.cdlStoreURL == nil {
		return ""
	}
	return c.cdlStoreURL.String()
}

// SetCDLStoreURL sets the CDL store URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetCDLStoreURL(urlStr string) error {
	if urlStr == "" {
		return ErrCDLURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.cdlStoreURL, err = url.Parse(urlStr)
	return err
}

func (c *Client) newCDLRequest(method, path string, opt interface{}, options ...OptionFunc) (*http.Request, error) {
	u := *c.cdlStoreURL
	// Set the encoded opaque data
	u.Opaque = c.cdlStoreURL.Path + path

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
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+c.iamClient.Token())
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
	Link []LinkElementType
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
