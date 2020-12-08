package console

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/fhir"
)

type tokenType int
type ContextKey string

const (
	libraryVersion    = "0.21.1"
	userAgent         = "go-hsdp-api/console/" + libraryVersion
	consoleAPIVersion = "3"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Jti          string `json:"jti"`
}

type CFLinksResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		CloudControllerV2 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v2"`
		CloudControllerV3 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v3"`
		NetworkPolicyV0 struct {
			Href string `json:"href"`
		} `json:"network_policy_v0"`
		NetworkPolicyV1 struct {
			Href string `json:"href"`
		} `json:"network_policy_v1"`
		Login struct {
			Href string `json:"href"`
		} `json:"login"`
		UAA struct {
			Href string `json:"href"`
		} `json:"uaa"`
		Credhub interface{} `json:"credhub"`
		Routing struct {
			Href string `json:"href"`
		} `json:"routing"`
		Logging struct {
			Href string `json:"href"`
		} `json:"logging"`
		LogCache struct {
			Href string `json:"href"`
		} `json:"log_cache"`
		LogStream struct {
			Href string `json:"href"`
		} `json:"log_stream"`
		AppSSH struct {
			Href string `json:"href"`
			Meta struct {
				HostKeyFingerprint string `json:"host_key_fingerprint"`
				OauthClient        string `json:"oauth_client"`
			} `json:"meta"`
		} `json:"app_ssh"`
	} `json:"links"`
}

const (
	oAuthToken tokenType = iota
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// A Client manages communication with HSDP IAM API
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	config *Config

	validate *validator.Validate

	baseConsoleURL *url.URL
	baseUAAURL     *url.URL

	// token type used to make authenticated API calls.
	tokenType tokenType

	// tokens used to make authenticated API calls.
	token        string
	refreshToken string
	idToken      string
	expiresAt    time.Time

	// scope holds the client scope
	scopes []string

	// User agent used when communicating with the HSDP IAM API.
	UserAgent string

	Metrics *MetricsService

	debugFile *os.File

	sync.Mutex
}

// NewClient returns a new HSDP Console API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide a valid oAuth bearer token.
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	return newClient(httpClient, config)
}

func newClient(httpClient *http.Client, config *Config) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if config.UAAURL == "" && config.BaseConsoleURL == "" {
		doAutoconf(config)
	}
	if config.UAAURL == "" {
		return nil, ErrUAAURLCannotBeEmpty
	}
	if config.BaseConsoleURL == "" {
		return nil, ErrConsoleURLCannotBeEmpty
	}
	c := &Client{client: httpClient, config: config, UserAgent: userAgent}
	if err := c.SetBaseUAAURL(c.config.UAAURL); err != nil {
		return nil, err
	}
	if err := c.SetBaseConsoleURL(c.config.BaseConsoleURL); err != nil {
		return nil, err
	}

	if config.DebugLog != "" {
		debugFile, err := os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		} else {
			c.debugFile = debugFile
		}
	}
	c.Metrics = &MetricsService{client: c}
	c.validate = validator.New()
	return c, nil
}

func doAutoconf(config *Config) {
	if config.Region != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region))
		if err == nil {
			uaaService := c.Service("uaa")
			consoleService := c.Service("console")
			if config.UAAURL == "" {
				config.UAAURL = uaaService.URL
			}
			if config.BaseConsoleURL == "" {
				config.BaseConsoleURL = consoleService.URL
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

// Returns the http Client used for connections
func (c *Client) HttpClient() *http.Client {
	return c.client
}

// Token returns the current token
func (c *Client) Token() string {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	expires := c.expiresAt.Unix()

	if expires-now < 60 {
		if c.TokenRefresh() != nil {
			return ""
		}
	}
	return c.token
}

// TokenRefresh refreshes the accessToken
func (c *Client) TokenRefresh() error {
	if c.refreshToken == "" {
		return ErrMissingRefreshToken
	}

	u := *c.baseUAAURL
	u.Opaque = c.baseUAAURL.Path + "oauth/token"

	req := &http.Request{
		Method:     "POST",
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	form := url.Values{}
	form.Add("token", c.refreshToken)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", c.refreshToken)
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth("cf", "")
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

// RefreshToken returns the refresh token
func (c *Client) RefreshToken() string {
	return c.refreshToken
}

// IDToken returns the ID token
func (c *Client) IDToken() string {
	return c.idToken
}

// Expires returns the expiry time (Unix) of the access token
func (c *Client) Expires() int64 {
	return c.expiresAt.Unix()
}

// SetBaseConsoleURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseConsoleURL(urlStr string) error {
	if urlStr == "" {
		return ErrConsoleURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseConsoleURL, err = url.Parse(urlStr)
	return err
}

// SetBaseIDMURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseUAAURL(urlStr string) error {
	if urlStr == "" {
		return ErrUAAURLCannotBeEmpty
	}
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseUAAURL, err = url.Parse(urlStr)
	return err
}

// Endpoint type
type Endpoint string

// Constants
const (
	UAA     = "UAA"
	CONSOLE = "CONSOLE"
)

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(endpoint, method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	var u url.URL
	switch endpoint {
	case UAA:
		u = *c.baseUAAURL
		u.Opaque = c.baseUAAURL.Path + path
	case CONSOLE:
		u = *c.baseConsoleURL
		u.Opaque = c.baseConsoleURL.Path + path
	default:
		return nil, fmt.Errorf("Unknown endpoint: `%s`", endpoint)
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

	switch c.tokenType {
	case oAuthToken:
		if token := c.Token(); token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response is a HSDP Console API response. This wraps the standard http.Response
// returned from HSDP Console and provides convenient access to things like errors
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	id := uuid.New()

	if c.config.Debug {
		dumped, _ := httputil.DumpRequest(req, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Request [%s] start ---\n%s\n[go-hsdp-api] --- Request [%s] end ---\n", id, string(dumped), id)
		if c.debugFile != nil {
			_, _ = c.debugFile.WriteString(out)
		} else {
			fmt.Println(out)
		}
	}
	resp, err := c.client.Do(req)
	if c.config.Debug && resp != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response [%s] start ---\n%s\n[go-hsdp-api] --- Response [%s] end ---\n", id, string(dumped), id)
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

// WithContext runs the request with the provided context
func WithContext(ctx context.Context) OptionFunc {
	return func(req *http.Request) error {
		*req = *req.WithContext(ctx)
		return nil
	}
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
