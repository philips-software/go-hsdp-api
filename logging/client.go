// Package logging provides support for HSDP Logging services
package logging

import (
	"bytes"
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

	"github.com/philips-software/go-hsdp-api/fhir"
	signer "github.com/philips-software/go-hsdp-signer"
	errors "golang.org/x/xerrors"
)

const (
	libraryVersion = "0.1.0"
	userAgent      = "go-hsdp-api/logging/" + libraryVersion
)

var (
	// LogTimeFormat is the log time format to use
	LogTimeFormat = "2006-01-02T15:04:05.000Z07:00"

	ErrNothingToPost       = errors.New("nothing to post")
	ErrMissingSharedKey    = errors.New("missing shared key")
	ErrMissingSharedSecret = errors.New("missing shared secret")
	ErrMissingBaseURL      = errors.New("missing base URL")
	ErrMissingProductKey   = errors.New("missing ProductKey")

	scaryMap = map[string]string{
		";":  "💀",
		"\\": "🎃",
		"&":  "👻",
		">":  "👿",
		"<":  "👾",
	}
)

// Storer defines the store operations for logging
type Storer interface {
	StoreResources(msgs []Resource, count int) (*Response, error)
}

// Config the client
type Config struct {
	SharedKey    string
	SharedSecret string
	BaseURL      string
	ProductKey   string
	Debug        bool
}

// Valid returns if all required config fields are present, false otherwise
func (c *Config) Valid() (bool, error) {
	if c.SharedKey == "" {
		return false, ErrMissingSharedKey
	}
	if c.SharedSecret == "" {
		return false, ErrMissingSharedSecret
	}
	if c.BaseURL == "" {
		return false, ErrMissingBaseURL
	}
	if c.ProductKey == "" {
		return false, ErrMissingProductKey
	}
	return true, nil
}

// Client holds the client state
type Client struct {
	config     Config
	url        *url.URL
	httpClient *http.Client
	httpSigner *signer.Signer
}

// Response holds a LogEvent response
type Response struct {
	*http.Response
	Message string
	Failed  []Resource
}

// CustomIndexBody describes the custom index request payload
type CustomIndexBody []struct {
	Fieldname string `json:"fieldname"`
	Fieldtype string `json:"fieldtype"`
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// NewClient returns an instance of the logger client with the given Config
func NewClient(httpClient *http.Client, config Config) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if valid, err := config.Valid(); !valid {
		return nil, err
	}
	var logger Client

	logger.config = config
	logger.httpClient = httpClient

	url, err := url.Parse(config.BaseURL + "/core/log/LogEvent")
	if err != nil {
		return nil, err
	}

	logger.httpSigner, err = signer.New(logger.config.SharedKey, logger.config.SharedSecret)
	if err != nil {
		return nil, err
	}

	logger.url = url
	if os.Getenv("DEBUG") == "true" {
		logger.config.Debug = true
	}
	return &logger, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	if c.config.Debug {
		dumped, _ := httputil.DumpRequest(req, true)
		fmt.Fprintf(os.Stderr, "REQUEST: %s\n", string(dumped))
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
			fmt.Fprintf(os.Stderr, "RESPONSE: %s\n", string(dumped))
		} else {
			fmt.Fprintf(os.Stderr, "Error sending response: %s\n", err)
		}
	}

	err = CheckResponse(resp)
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

// ErrorResponse holds an error response from the server
type ErrorResponse struct {
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Opaque)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = "failed to parse unknown error format"
		}

		errorResponse.Message = fhir.ParseError(raw)
	}

	return errorResponse
}

// StoreResources posts one or more log messages
func (c *Client) StoreResources(msgs []Resource, count int) (*Response, error) {
	var b Bundle
	var invalid []Resource

	b.ResourceType = "Bundle"
	b.Entry = make([]Element, count)
	b.Type = "transaction"
	b.ProductKey = c.config.ProductKey

	j := 0
	for i := 0; i < count; i++ {
		msg := msgs[i]
		replaceScaryCharacters(&msg)
		if !msg.Valid() {
			invalid = append(invalid, msg)
			continue
		}
		// Element
		var e Element
		e.Resource = msg
		// Base64 encode here
		e.Resource.LogData.Message = base64.StdEncoding.EncodeToString([]byte(msg.LogData.Message))
		e.Resource.ResourceType = "LogEvent"
		b.Entry[j] = e
		j++
	}
	if j == 0 { // No payload
		return nil, ErrNothingToPost
	}

	b.Total = j

	req := &http.Request{
		Method:     http.MethodPost,
		URL:        c.url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       c.url.Host,
	}

	bodyBytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyBytes)
	req.Body = ioutil.NopCloser(bodyReader)
	req.ContentLength = int64(bodyReader.Len())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Version", "1")
	req.Header.Set("User-Agent", userAgent)
	if err := c.httpSigner.SignRequest(req); err != nil {
		return nil, err
	}

	var serverResponse bytes.Buffer

	resp, err := c.Do(req, &serverResponse)

	if len(invalid) > 0 {
		resp.Failed = invalid
	}

	return resp, err
}

func replaceScaryCharacters(msg *Resource) {
	if len(msg.Custom) == 0 {
		return
	}
	stringCustom := string(msg.Custom)
	for s, r := range scaryMap {
		stringCustom = strings.ReplaceAll(stringCustom, s, r)
	}
	msg.Custom = []byte(stringCustom)
}
