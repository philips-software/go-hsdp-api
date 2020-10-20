// Package cartel provides support for HSDP Cartel services
package cartel

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	autoconf "github.com/philips-software/go-hsdp-api/config"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.21.1"
	userAgent      = "go-hsdp-api/cartel/" + libraryVersion
)

// Config the client
type Config struct {
	Region     string `cloud:"-" json:"-"`
	Token      string `cloud:"token" json:"token"`
	Secret     string `cloud:"secret" json:"secret"`
	SkipVerify bool   `cloud:"skip_verify" json:"skip_verify"`
	NoTLS      bool   `cloud:"no_tls" json:"no_tls"`
	Host       string `cloud:"host" json:"host"`
	Debug      bool   `cloud:"-" json:"debug,omitempty"`
	DebugLog   string `cloud:"-" json:"debug_log,omitempty"`
}

// Valid returns if all required config fields are present, false otherwise
func (c *Config) Valid() (bool, error) {
	if len(c.Secret) == 0 {
		return false, ErrMissingSecret
	}
	if len(c.Token) == 0 {
		return false, ErrMissingToken
	}
	if c.Host == "" {
		return false, ErrMissingHost
	}
	return true, nil
}

// Client holds the client state
type Client struct {
	config     *Config
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string

	debugFile *os.File
}

// Response holds a LogEvent response
type Response struct {
	*http.Response
	Message string
}

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func doAutoconf(config *Config) {
	if config.Region != "" {
		ac, err := autoconf.New(
			autoconf.WithRegion(config.Region))
		if err == nil {
			loggingService := ac.Service("cartel")
			if host, err := loggingService.GetString("host"); err == nil && config.Host == "" {
				config.Host = host
			}
		}
	}
}

// NewClient returns an instance of the logger client with the given Config
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
		tr := &http.Transport{}
		if config.SkipVerify {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			tr.TLSClientConfig = &tls.Config{}
		}
		httpClient.Transport = tr
	}
	doAutoconf(config)
	if valid, err := config.Valid(); !valid {
		return nil, err
	}
	var cartel Client

	cartel.config = config
	cartel.httpClient = httpClient
	cartel.userAgent = userAgent

	// Make sure the given URL ends with a slash
	host := fmt.Sprintf("https://%s", cartel.config.Host)
	if config.NoTLS {
		host = fmt.Sprintf("http://%s", cartel.config.Host)
	}
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	var err error
	cartel.baseURL, err = url.Parse(host)
	if err != nil {
		return nil, err
	}

	configDebug(&cartel)
	return &cartel, nil
}

func configDebug(cartel *Client) {
	if cartel.config.Debug {
		if cartel.config.DebugLog != "" {
			debugFile, err := os.OpenFile(cartel.config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err == nil {
				cartel.debugFile = debugFile
			}
		}
		if cartel.debugFile == nil {
			cartel.debugFile = os.Stderr
		}
	}
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	if c.config.Debug {
		dumped, _ := httputil.DumpRequest(req, true)
		fmt.Fprintf(c.debugFile, "REQUEST: %s\n", string(dumped))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	if c.config.Debug {
		if resp != nil {
			dumped, _ := httputil.DumpResponse(resp, true)
			_, _ = fmt.Fprintf(c.debugFile, "RESPONSE: %s\n", string(dumped))
		} else {
			_, _ = fmt.Fprintf(c.debugFile, "Error sending response: %s\n", err)
		}
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	if err != nil {
		return response, err
	}
	err = CheckResponse(resp)
	return response, err
}

// ErrorResponse holds an error response from the server
type ErrorResponse struct {
	Response    *http.Response `json:"-"`
	Message     string         `json:"-"`
	Code        int            `json:"code"`
	Description string         `json:"description"`
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}
	return ErrNonHttp20xResponse
}

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, opt *RequestBody, options []OptionFunc) (*http.Request, error) {
	u := *c.baseURL
	u.Opaque = c.baseURL.Path + path

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
	// Add token
	if opt != nil && opt.Token == "" {
		opt.Token = c.config.Token
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
		req.Header.Set("Authorization", string(generateSignature([]byte(c.config.Secret), bodyBytes)))
	}
	req.Header.Set("Accept", "application/json")

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

func generateSignature(secret []byte, payload []byte) string {
	hash := hmac.New(sha256.New, secret)
	_, _ = hash.Write(payload)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
