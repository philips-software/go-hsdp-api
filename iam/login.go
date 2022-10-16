package iam

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CodeLogin uses the authorization_code grant type to fetch tokens
func (c *Client) CodeLogin(code string, redirectURI string) error {
	// Authorize
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
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	if len(redirectURI) > 0 {
		form.Add("redirect_uri", redirectURI)
	}
	body := form.Encode()
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = io.NopCloser(strings.NewReader(body))
	req.ContentLength = int64(len(body))

	return c.doTokenRequest(req)
}

// ServiceLogin logs a service in using a JWT signed with the service private key
func (c *Client) ServiceLogin(service Service) error {
	accessTokenEndpoint := c.accessTokenEndpoint()
	token, err := service.GenerateJWT(accessTokenEndpoint)
	if err != nil {
		return err
	}
	// Authorize
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
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	// HSDP IAM currently croaks on URL encoded grant_type value. INC0038532
	body := "assertion=" + token
	body += "&grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer"
	body += "&"
	body += form.Encode()

	// Remove Authorization header
	req.Header.Del("Authorization")

	req.Body = io.NopCloser(strings.NewReader(body))
	req.ContentLength = int64(len(body))
	c.service = service // Save service so we can refresh later!

	return c.doTokenRequest(req)
}

// Login logs in a user with `username` and `password`
func (c *Client) Login(username, password string) error {
	// Authorize
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
	form.Add("username", username)
	form.Add("password", password)
	form.Add("grant_type", "password")
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))
	c.service = Service{} // reset

	return c.doTokenRequest(req)
}

// ClientCredentialsLogin logs in using client credentials
// The client credentials and scopes are expected to passed during configuration of the client
func (c *Client) ClientCredentialsLogin() error {
	// Authorize
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
	form.Add("grant_type", "client_credentials")
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

// RevokeAccessToken revokes the access and refresh token
func (c *Client) RevokeAccessToken() error {
	return c.revokeToken(c.token)
}

// RevokeRefreshAccessToken revokes the access and refresh token
func (c *Client) RevokeRefreshAccessToken() error {
	return c.revokeToken(c.refreshToken)
}

type endSessionOptions struct {
	IDTokenHint *string `url:"id_token_hint,omitempty"`
}

// EndSession ends the current active session
func (c *Client) EndSession() error {
	req, err := c.newRequest(IAM, "GET", "authorize/oauth2/endsession", &endSessionOptions{
		IDTokenHint: &c.idToken,
	}, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

func (c *Client) revokeToken(token string) error {
	req, err := c.newRequest(IAM, "POST", "authorize/oauth2/revoke", nil, nil)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("token", token)
	if len(c.config.Scopes) > 0 {
		scopes := strings.Join(c.config.Scopes, " ")
		form.Add("scope", scopes)
	}
	req.SetBasicAuth(c.config.OAuth2ClientID, c.config.OAuth2Secret)
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
}

func (c *Client) doTokenRequest(req *http.Request) error {
	var tokenResponse tokenResponse

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Api-Version", loginAPIVersion)
	resp, err := c.do(req, &tokenResponse)
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("login failed: %d", resp.StatusCode())
	}
	if tokenResponse.AccessToken == "" {
		return ErrNotAuthorized
	}
	c.tokenType = OAuthToken
	c.token = tokenResponse.AccessToken
	if tokenResponse.RefreshToken != "" { // Doesn't always contain new refresh token
		c.refreshToken = tokenResponse.RefreshToken
	}
	if tokenResponse.IDToken != "" {
		c.idToken = tokenResponse.IDToken
	}
	c.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	c.scopes = strings.Split(tokenResponse.Scope, " ")
	return nil
}
