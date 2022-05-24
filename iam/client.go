// Package iam provides support for interacting with HSDP IAM and IDM services
package iam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/philips-software/go-hsdp-api/internal"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	hsdpsigner "github.com/philips-software/go-hsdp-signer"
)

type tokenType int
type ContextKey string

const (
	userAgent       = "go-hsdp-api/iam/" + internal.LibraryVersion
	loginAPIVersion = "2"
)

type tokenResponse struct {
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

const (
	OAuthToken tokenType = iota
	JWTToken   tokenType = 1
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// A Client manages communication with HSDP IAM API
type Client struct {
	// HTTP client used to communicate with the API.
	*http.Client

	config *Config

	signer   *hsdpsigner.Signer
	validate *validator.Validate

	baseIAMURL *url.URL
	baseIDMURL *url.URL

	// token type used to make authenticated API calls.
	tokenType tokenType

	// tokens used to make authenticated API calls.
	token        string
	refreshToken string
	idToken      string
	expiresAt    time.Time
	service      Service

	// scope holds the client scope
	scopes []string

	// User agent used when communicating with the HSDP IAM API.
	UserAgent string

	debugFile *os.File

	Organizations    *OrganizationsService
	Groups           *GroupsService
	Permissions      *PermissionsService
	Roles            *RolesService
	Users            *UsersService
	Applications     *ApplicationsService
	Propositions     *PropositionsService
	Clients          *ClientsService
	Services         *ServicesService
	MFAPolicies      *MFAPoliciesService
	PasswordPolicies *PasswordPoliciesService
	Devices          *DevicesService
	EmailTemplates   *EmailTemplatesService
	SMSGateways      *SMSGatewaysService
	SMSTemplates     *SMSTemplatesService

	sync.Mutex
}

// NewClient returns a new HSDP IAM API client. If a nil httpClient is
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
	doAutoconf(config)
	c := &Client{Client: httpClient, config: config, UserAgent: userAgent}
	if err := c.SetBaseIAMURL(c.config.IAMURL); err != nil {
		return nil, err
	}
	if err := c.SetBaseIDMURL(c.config.IDMURL); err != nil {
		return nil, err
	}
	if config.Signer == nil {
		signer, err := hsdpsigner.New(c.config.SharedKey, c.config.SecretKey)
		if err != nil { // Allow nil signer
			signer = nil
		}
		c.signer = signer
	} else {
		c.signer = config.Signer
	}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			httpClient.Transport = internal.NewLoggingRoundTripper(httpClient.Transport, c.debugFile)
		}
	}

	c.validate = validator.New()
	c.Organizations = &OrganizationsService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Permissions = &PermissionsService{client: c}
	c.Roles = &RolesService{client: c}
	c.Users = &UsersService{client: c, validate: validator.New()}
	c.Applications = &ApplicationsService{client: c}
	c.Propositions = &PropositionsService{client: c}
	c.Clients = &ClientsService{client: c, validate: validator.New()}
	c.Services = &ServicesService{client: c}
	c.MFAPolicies = &MFAPoliciesService{client: c, validate: validator.New()}
	c.PasswordPolicies = &PasswordPoliciesService{client: c, validate: validator.New()}
	c.Devices = &DevicesService{client: c, validate: validator.New()}
	c.EmailTemplates = &EmailTemplatesService{client: c, validate: validator.New()}
	c.SMSGateways = &SMSGatewaysService{client: c, validate: validator.New()}
	c.SMSTemplates = &SMSTemplatesService{client: c, validate: validator.New()}
	return c, nil
}

func doAutoconf(config *Config) {
	if config.Region != "" && config.Environment != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region),
			autoconf.WithEnv(config.Environment))
		if err == nil {
			iamService := c.Service("iam")
			idmService := c.Service("idm")
			if iamService.URL != "" && config.IAMURL == "" {
				config.IAMURL = iamService.URL
			}
			if idmService.URL != "" && config.IDMURL == "" {
				config.IDMURL = idmService.URL
			}
		}
	}
}

func (c *Client) validSigner() bool {
	return c.signer != nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

// HttpClient returns the http Client used for connections
func (c *Client) HttpClient() *http.Client {
	return c.Client
}

// WithToken returns a cloned client with the token set
func (c *Client) WithToken(token string) *Client {
	client, _ := NewClient(c.Client, c.config)
	client.SetToken(token)
	return client
}

// WithLogin returns a cloned client with new login
func (c *Client) WithLogin(username, password string) (*Client, error) {
	client, err := NewClient(c.Client, c.config)
	if err != nil {
		return nil, err
	}
	err = client.Login(username, password)
	return client, err
}

func (c *Client) accessTokenEndpoint() string {
	if c.baseIAMURL != nil {
		return c.baseIAMURL.String() + "oauth2/access_token"
	}
	return ""
}

// Token returns the current token
func (c *Client) Token() (string, error) {
	now := time.Now().Unix()
	expires := c.expiresAt.Unix()

	if expires-now < 60 {
		if err := c.TokenRefresh(); err != nil {
			return "", err
		}
	}
	c.Lock()
	defer c.Unlock()
	return c.token, nil
}

// ExpireToken expires the token immediately
func (c *Client) ExpireToken() {
	c.Lock()
	defer c.Unlock()
	c.expiresAt = time.Now()
}

// TokenRefresh forces a token refresh
func (c *Client) TokenRefresh() error {
	c.Lock()
	defer c.Unlock()

	if c.refreshToken == "" {
		if c.service.Valid() { // Possible service
			return c.ServiceLogin(c.service)
		}
		return ErrMissingRefreshToken
	}

	u := *c.baseIAMURL
	u.Opaque = c.baseIAMURL.Path + "authorize/oauth2/token"

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
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", c.refreshToken)
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	if !c.HasOAuth2Credentials() {
		return ErrMissingOAuth2Credentials
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

// HasOAuth2Credentials returns true if the client is configured with OAuth2 credentials
func (c *Client) HasOAuth2Credentials() bool {
	return c.config.OAuth2ClientID != "" && c.config.OAuth2Secret != ""
}

// HasScopes returns true of all scopes are there for the client
func (c *Client) HasScopes(scopes ...string) bool {
	for _, s := range scopes {
		found := false
		for _, t := range c.scopes {
			if t == s {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// HasPermissions returns true if all permissions are there for the client
func (c *Client) HasPermissions(orgID string, permissions ...string) bool {
	introspect, _, err := c.Introspect(WithOrgContext(orgID))
	if err != nil {
		return false
	}
	foundOrg := false
	for _, org := range introspect.Organizations.OrganizationList {
		if org.OrganizationID != orgID {
			continue
		}
		foundOrg = true
		// Search in the organization effective permissions list
		for _, p := range permissions {
			found := false
			for _, q := range org.EffectivePermissions {
				if p == q {
					found = true
					continue
				}
			}
			if !found {
				// Permission is missing so return false
				return false
			}
		}
	}
	return foundOrg
}

// SetToken sets the token
func (c *Client) SetToken(token string) {
	c.token = token
	c.expiresAt = time.Now().Add(86400 * time.Second)
	c.tokenType = OAuthToken
}

// SetTokens sets the token
func (c *Client) SetTokens(accessToken, refreshToken, idToken string, expiresAt int64) {
	c.Lock()
	defer c.Unlock()
	c.token = accessToken
	c.refreshToken = refreshToken
	c.idToken = idToken
	c.expiresAt = time.Unix(expiresAt, 0)
	c.tokenType = OAuthToken
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

// BaseIAMURL return a copy of the baseIAMURL.
func (c *Client) BaseIAMURL() *url.URL {
	u := *c.baseIAMURL
	return &u
}

// BaseIDMURL return a copy of the baseIAMURL.
func (c *Client) BaseIDMURL() *url.URL {
	u := *c.baseIDMURL
	return &u
}

// SetBaseIAMURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseIAMURL(urlStr string) error {
	if urlStr == "" {
		return ErrBaseIAMCannotBeEmpty
	}
	// Make sure the given URL ends with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseIAMURL, err = url.Parse(urlStr)
	return err
}

// SetBaseIDMURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseIDMURL(urlStr string) error {
	if urlStr == "" {
		return ErrBaseIDMCannotBeEmpty
	}
	// Make sure the given URL ends with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseIDMURL, err = url.Parse(urlStr)
	return err
}

// Endpoint type
type Endpoint string

// Constants
const (
	IAM = "IAM"
	IDM = "IDM"
)

// newRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newRequest(endpoint, method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	var u url.URL
	switch endpoint {
	case IDM:
		u = *c.baseIDMURL
		u.Opaque = c.baseIDMURL.Path + path
	case IAM:
		u = *c.baseIAMURL
		u.Opaque = c.baseIAMURL.Path + path
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
	req.Header.Set("User-Agent", userAgent)

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
	case OAuthToken:
		if token, err := c.Token(); err == nil {
			req.Header.Set("Authorization", "Bearer "+token)
		} else {
			req.Header.Set("X-Token-Error", fmt.Sprintf("%v", err))
		}
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

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// doSigned performs a signed API request
func (c *Client) doSigned(req *http.Request, v interface{}) (*Response, error) {
	if c.signer == nil {
		return nil, ErrNoValidSignerAvailable
	}
	if err := c.signer.SignRequest(req); err != nil {
		return nil, err
	}
	return c.do(req, v)
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	response := newResponse(resp)

	err = internal.CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil && response.StatusCode != http.StatusNoContent {
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

// ErrorResponse represents an IAM errors response
// containing a code and a human-readable message
type ErrorResponse struct {
	Response         *http.Response `json:"-"`
	Code             string         `json:"responseCode,omitempty"`
	Message          string         `json:"responseMessage,omitempty"`
	ErrorString      string         `json:"error,omitempty"`
	ErrorDescription string         `json:"error_description,omitempty"`
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Opaque)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}
