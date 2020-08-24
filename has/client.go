// Package has provides support for HSDP Appstreaming service
package has

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
	libraryVersion = "0.21.0"
	userAgent      = "go-hsdp-api/has/" + libraryVersion
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	HASURL   string
	OrgID    string
	Debug    bool
	DebugLog string
}

// A Client manages communication with HSDP IAM API
type Client struct {
	// HTTP client used to communicate with the API.
	iamClient *iam.Client

	config *Config

	baseHASURL *url.URL

	// User agent used when communicating with the HSDP IAM API.
	UserAgent string

	debugFile *os.File

	Resources *ResourcesService
	Sessions  *SessionsService
	Images    *ImagesService
}

// NewClient returns a new HSDP HAS API client. If a nil httpClient is
// provided, http.DefaultClient will be used. A configured IAM client must be provided
// as well
func NewClient(iamClient *iam.Client, config *Config) (*Client, error) {
	return newClient(iamClient, config)
}

func newClient(iamClient *iam.Client, config *Config) (*Client, error) {
	c := &Client{iamClient: iamClient, config: config, UserAgent: userAgent}
	if err := c.SetBaseHASURL(c.config.HASURL); err != nil {
		return nil, err
	}
	if config.OrgID == "" || !iamClient.HasPermissions(config.OrgID,
		"HAS_SESSION.ALL", "HAS_RESOURCE.ALL") {
		return nil, ErrMissingHASPermissions
	}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		}
	}

	c.Resources = &ResourcesService{client: c, orgID: config.OrgID}
	c.Sessions = &SessionsService{client: c, orgID: config.OrgID}
	c.Images = &ImagesService{client: c, orgID: config.OrgID}
	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		c.debugFile.Close()
		c.debugFile = nil
	}
}

// SetBaseHASURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseHASURL(urlStr string) error {
	if urlStr == "" {
		return ErrBaseHASCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseHASURL, err = url.Parse(urlStr)
	return err
}

// NewHASRequest creates an new HAS API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewHASRequest(method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	u := *c.baseHASURL
	// Set the encoded opaque data
	u.Opaque = c.baseHASURL.Path + path

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
			fmt.Print(out)
		}
	}
	resp, err := c.iamClient.HttpClient().Do(req)
	if c.config.Debug && resp != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response start ---\n%s\n[go-hsdp-api] --- Response end ---\n", string(dumped))
		if c.debugFile != nil {
			_, _ = c.debugFile.WriteString(out)
		} else {
			fmt.Print(out)
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
