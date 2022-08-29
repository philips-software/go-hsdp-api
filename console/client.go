// Package console provides support for HSDP Console APIs
package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/philips-software/go-hsdp-api/internal"
	"golang.org/x/oauth2"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	autoconf "github.com/philips-software/go-hsdp-api/config"
)

type tokenType int
type ContextKey string

const (
	userAgent = "go-hsdp-api/console/" + internal.LibraryVersion
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
	*http.Client

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

	debugFile  *os.File
	consoleErr error

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
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	if config.UAAURL == "" && config.BaseConsoleURL == "" {
		doAutoconf(config)
	}
	if config.UAAURL == "" {
		return nil, ErrUAAURLCannotBeEmpty
	}
	c := &Client{Client: httpClient, config: config, UserAgent: userAgent}
	if err := c.SetBaseUAAURL(c.config.UAAURL); err != nil {
		return nil, err
	}
	if err := c.SetBaseConsoleURL(c.config.BaseConsoleURL); err != nil {
		c.consoleErr = err
	}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			httpClient.Transport = internal.NewLoggingRoundTripper(httpClient.Transport, c.debugFile)
		}
	}
	header := make(http.Header)
	header.Set("User-Agent", userAgent)
	httpClient.Transport = internal.NewHeaderRoundTripper(httpClient.Transport, header)

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
	return c.Client
}

// Token returns the current token. It also conforms to TokenSource
func (c *Client) Token() (*oauth2.Token, error) {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	expires := c.expiresAt.Unix()

	if expires-now < 60 {
		if c.TokenRefresh() != nil {
			return nil, fmt.Errorf("failed to refresh console token")
		}
	}
	return &oauth2.Token{
		AccessToken:  c.token,
		RefreshToken: c.refreshToken,
		Expiry:       c.expiresAt,
	}, nil
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
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
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

func (c *Client) newRequest(endpoint, method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	var u url.URL
	switch endpoint {
	case UAA:
		u = *c.baseUAAURL
		u.Opaque = c.baseUAAURL.Path + path
	case CONSOLE:
		if c.consoleErr != nil {
			return nil, c.consoleErr
		}
		u = *c.baseConsoleURL
		u.Opaque = c.baseConsoleURL.Path + path
	default:
		return nil, fmt.Errorf("unknown endpoint: `%s`", endpoint)
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
		req.Body = io.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	switch c.tokenType {
	case oAuthToken:
		if token, err := c.Token(); err == nil {
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		}
	}
	return req, nil
}

// Response is a HSDP Console API response. This wraps the standard http.Response
// returned from HSDP Console and provides convenient access to things like errors
type Response struct {
	*http.Response
	Error
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return response, err
		}
	}
	err = internal.CheckResponse(resp)
	return response, err
}
