// Package iam provides support for interacting with HSDP IAM/IDM services
package iam

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
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/philips-software/go-hsdp-api/fhir"
	hsdpsigner "github.com/philips-software/go-hsdp-signer"
)

type tokenType int
type ContextKey string

const (
	libraryVersion                 = "0.19.0"
	userAgent                      = "go-hsdp-api/iam/" + libraryVersion
	loginAPIVersion                = "2"
	ContextKeyRequestID ContextKey = "requestID"
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
	oAuthToken tokenType = iota
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	OAuth2ClientID   string
	OAuth2Secret     string
	SharedKey        string
	SecretKey        string
	BaseIAMURL       string
	BaseIDMURL       string
	OrgAdminUsername string
	OrgAdminPassword string
	IAMURL           string
	IDMURL           string
	Scopes           []string
	RootOrgID        string
	Debug            bool
	DebugLog         string
	Signer           *hsdpsigner.Signer
}

// A Client manages communication with HSDP IAM API
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

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
	expiresAt    time.Time

	// scope holds the client scope
	scopes []string

	// User agent used when communicating with the HSDP IAM API.
	UserAgent string

	debugFile *os.File

	Organizations *OrganizationsService
	Groups        *GroupsService
	Permissions   *PermissionsService
	Roles         *RolesService
	Users         *UsersService
	Applications  *ApplicationsService
	Propositions  *PropositionsService
	Clients       *ClientsService
	Services      *ServicesService
	MFAPolicies   *MFAPoliciesService
}

// NewClient returns a new HSDP IAM API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide a valid oAuth bearer token.
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	return newClient(httpClient, config)
}

func newClient(httpClient *http.Client, config *Config) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient, config: config, UserAgent: userAgent}
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
		debugFile, err := os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		} else {
			c.debugFile = debugFile
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
	return c, nil
}

func (c *Client) validSigner() bool {
	return c.signer != nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		c.debugFile.Close()
		c.debugFile = nil
	}
}

// Returns the http Client used for connections
func (c *Client) HttpClient() *http.Client {
	return c.client
}

// WithToken returns a cloned client with the token set
func (c *Client) WithToken(token string) *Client {
	client, _ := NewClient(c.client, c.config)
	client.SetToken(token)
	return client
}

// WithLogin returns a cloned client with new login
func (c *Client) WithLogin(username, password string) (*Client, error) {
	client, err := NewClient(c.client, c.config)
	if err != nil {
		return nil, err
	}
	err = client.Login(username, password)
	return client, err
}

func (c *Client) accessTokenEndpoint() string {
	return c.baseIAMURL.String() + "oauth2/access_token"
}

// CodeLogin
func (c *Client) CodeLogin(code string, redirectURI string) error {
	// Authorize
	req, err := c.NewRequest(IAM, "POST", "authorize/oauth2/token", nil, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	if len(redirectURI) > 0 {
		form.Add("redirect_uri", redirectURI)
	}
	body := form.Encode()
	req.Body = ioutil.NopCloser(strings.NewReader(body))
	req.ContentLength = int64(len(body))

	return c.doTokenRequest(req)
}

// ServiceLogin logs a service in using a JWT signed with the service private key
func (c *Client) ServiceLogin(service Service) error {
	token, err := service.GetToken(c.accessTokenEndpoint())
	if err != nil {
		return err
	}
	// Authorize
	req, err := c.NewRequest(IAM, "POST", "authorize/oauth2/token", nil, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	// HSDP IAM currently croakes on URL encoded grant_type value. INC0038532
	body := "assertion=" + token
	body += "&grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer"
	body += "&"
	body += form.Encode()

	req.Body = ioutil.NopCloser(strings.NewReader(body))
	req.ContentLength = int64(len(body))

	return c.doTokenRequest(req)
}

// Login logs in a user with `username` and `password`
func (c *Client) Login(username, password string) error {
	req, err := c.NewRequest(IAM, "POST", "authorize/oauth2/token", nil, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("grant_type", "password")
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

func (c *Client) doTokenRequest(req *http.Request) error {
	var tokenResponse tokenResponse

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Api-Version", loginAPIVersion)
	resp, err := c.Do(req, &tokenResponse)

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %d", resp.StatusCode)
	}
	if tokenResponse.AccessToken == "" {
		return ErrNotAuthorized
	}
	c.tokenType = oAuthToken
	c.token = tokenResponse.AccessToken
	c.refreshToken = tokenResponse.RefreshToken
	c.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	c.scopes = strings.Split(tokenResponse.Scope, " ")
	return nil
}

// Token returns the current token
func (c *Client) Token() string {
	now := time.Now().Unix()
	expires := c.expiresAt.Unix()

	if expires-now < 60 {
		if c.TokenRefresh() != nil {
			return ""
		}
	}
	return c.token
}

// TokenRefresh refreshes the current access token using the refresh token
func (c *Client) TokenRefresh() error {
	if c.refreshToken == "" {
		return ErrMissingRefreshToken
	}

	req, err := c.NewRequest(IAM, "POST", "authorize/oauth2/token", nil, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("token", c.refreshToken)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", c.refreshToken)
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
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
	intr, _, err := c.Introspect()
	if err != nil {
		return false
	}
	foundOrg := false
	for _, org := range intr.Organizations.OrganizationList {
		if org.OrganizationID != orgID {
			continue
		}
		foundOrg = true
		// Search in the organization permission list
		for _, p := range permissions {
			found := false
			for _, q := range org.Permissions {
				if p == q {
					found = true
					continue
				}
			}
			if !found {
				// Permission is missing to return false
				return false
			}
		}
	}
	return foundOrg
}

// SetToken sets the token
func (c *Client) SetToken(token string) {
	c.token = token
	c.expiresAt = time.Now().Add(86400 * time.Minute)
	c.tokenType = oAuthToken
}

// RefreshToken returns the refresh token
func (c *Client) RefreshToken() string {
	return c.refreshToken
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
	// Make sure the given URL end with a slash
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
	// Make sure the given URL end with a slash
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

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(endpoint, method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	var u url.URL
	switch endpoint {
	case IDM:
		u = *c.baseIDMURL
		u.Opaque = c.baseIDMURL.Path + path
	case IAM:
		u = *c.baseIAMURL
		u.Opaque = c.baseIAMURL.Path + path
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

// DoSigned performs a signed API request
func (c *Client) DoSigned(req *http.Request, v interface{}) (*Response, error) {
	if c.signer == nil {
		return nil, ErrNoValidSignerAvailable
	}
	if err := c.signer.SignRequest(req); err != nil {
		return nil, err
	}
	return c.Do(req, v)
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
