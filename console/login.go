package console

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Login logs in a user with `username` and `password`
func (c *Client) Login(username, password string) error {
	req, err := c.newRequest(UAA, "POST", "oauth/token", nil, nil)
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
	req.SetBasicAuth("cf", "")
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
	req.ContentLength = int64(len(form.Encode()))

	return c.doTokenRequest(req)
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

// SetToken sets the UAA token
func (c *Client) SetToken(token string) *Client {
	c.Lock()
	defer c.Unlock()

	c.token = token
	c.expiresAt = time.Now().Add(600 * time.Second)
	return c
}

// SetTokens sets the tokens
func (c *Client) SetTokens(accessToken, refreshToken, idToken string, expiresAt int64) {
	c.Lock()
	defer c.Unlock()
	c.token = accessToken
	c.refreshToken = refreshToken
	c.idToken = idToken
	c.expiresAt = time.Unix(expiresAt, 0)
	c.tokenType = oAuthToken
}

func (c *Client) doTokenRequest(req *http.Request) error {
	var tokenResponse tokenResponse

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.do(req, &tokenResponse)

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

func (c *Client) UserID() (string, error) {
	token, _ := jwt.Parse(c.IDToken(), func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == "" {
		return "", fmt.Errorf("invalid claims")
	}
	return claims["sub"].(string), nil
}
