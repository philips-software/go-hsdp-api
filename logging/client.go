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
	"regexp"
	"strconv"
	"strings"

	"go.elastic.co/apm/module/apmhttp"

	autoconf "github.com/philips-software/go-hsdp-api/config"

	signer "github.com/philips-software/go-hsdp-signer"
	errors "golang.org/x/xerrors"
)

const (
	// TimeFormat is the time format used for the LogTime field
	TimeFormat = "2006-01-02T15:04:05.000Z07:00"

	libraryVersion = "0.21.1"
	userAgent      = "go-hsdp-api/logging/" + libraryVersion
)

var (
	ErrNothingToPost       = errors.New("nothing to post")
	ErrMissingSharedKey    = errors.New("missing shared key")
	ErrMissingSharedSecret = errors.New("missing shared secret")
	ErrMissingBaseURL      = errors.New("missing base URL")
	ErrMissingProductKey   = errors.New("missing ProductKey")
	ErrBatchErrors         = errors.New("batch errors. check Invalid map for details")
	ErrResponseError       = errors.New("unexpected HSDP response error")

	entryRegex = regexp.MustCompile(`^entry\[(\d+)]`)

	scaryMap = map[string]string{
		";":    "[sc]",
		"&":    "[amp]",
		">":    "[gt]",
		"<":    "[lt]",
		"\\u":  "[utf]",
		"\\f":  "[ff]",
		"\\\"": "[qt]",
	}
)

// Config the client
type Config struct {
	Region       string
	Environment  string
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
	config     *Config
	url        *url.URL
	httpClient *http.Client
	httpSigner *signer.Signer
}

// Response holds a LogEvent response
type StoreResponse struct {
	*http.Response
	Message string
	Failed  map[int]Resource
}

// CustomIndexBody describes the custom index request payload
type CustomIndexBody []struct {
	Fieldname string `json:"fieldname"`
	Fieldtype string `json:"fieldtype"`
}

// NewClient returns an instance of the logger client with the given Config
func NewClient(httpClient *http.Client, config *Config) (*Client, error) {
	if httpClient == nil {
		httpClient = apmhttp.WrapClient(http.DefaultClient)
	}
	// Autoconfig
	if config.Region != "" && config.Environment != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region),
			autoconf.WithEnv(config.Environment))
		if err == nil {
			loggingService := c.Service("logging")
			if ingestorURL, err := loggingService.GetString("url" +
				""); err == nil && config.BaseURL == "" {
				config.BaseURL = ingestorURL
			}
		}
	}
	if valid, err := config.Valid(); !valid {
		return nil, err
	}
	var logger Client

	logger.config = config
	logger.httpClient = httpClient

	parsedURL, err := url.Parse(config.BaseURL + "/core/log/LogEvent")
	if err != nil {
		return nil, err
	}

	logger.httpSigner, err = signer.New(logger.config.SharedKey, logger.config.SharedSecret)
	if err != nil {
		return nil, err
	}

	logger.url = parsedURL
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
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.config.Debug {
		dumped, _ := httputil.DumpRequest(req, true)
		_, _ = fmt.Fprintf(os.Stderr, "REQUEST: %s\n", string(dumped))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.config.Debug {
		if resp != nil {
			dumped, _ := httputil.DumpResponse(resp, true)
			_, _ = fmt.Fprintf(os.Stderr, "RESPONSE: %s\n", string(dumped))
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Error sending response: %s\n", err)
		}
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
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

// StoreResources posts one or more log messages
// In case invalid resources are detected StoreResources will return
// with ErrBatchErrors and the Response.Failed map will contain the resources
// This also happens in case the HSDP Ingestor API flags resources. In both cases
// the complete batch should be considered as not persisted and the LogEvents should
// be resubmitted for storage
func (c *Client) StoreResources(msgs []Resource, count int) (*StoreResponse, error) {
	var b Bundle
	invalid := make(map[int]Resource)

	b.ResourceType = "Bundle"
	b.Entry = make([]Element, count)
	b.Type = "transaction"
	b.ProductKey = c.config.ProductKey

	j := 0
	for i := 0; i < count; i++ {
		msg := msgs[i]
		replaceScaryCharacters(&msg)
		if !msg.Valid() {
			invalid[i] = msg
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
	if len(invalid) > 0 { // Don't even POST anything due to errors in the batch
		resp := StoreResponse{
			Failed: invalid,
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
			},
		}
		return &resp, ErrBatchErrors
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

	return c.performAndParseResponse(req, msgs)
}

func (c *Client) performAndParseResponse(req *http.Request, msgs []Resource) (*StoreResponse, error) {
	invalid := make(map[int]Resource)

	var serverResponse bytes.Buffer

	resp, err := c.Do(req, &serverResponse)
	if resp == nil {
		return nil, err
	}
	storeResp := &StoreResponse{Response: resp}
	if resp.StatusCode != http.StatusCreated { // Only good outcome
		var errResponse bundleErrorResponse
		err := json.Unmarshal(serverResponse.Bytes(), &errResponse)
		if err != nil {
			return storeResp, err
		}
		if len(errResponse.Issue) == 0 || len(errResponse.Issue[0].Location) == 0 {
			return storeResp, ErrResponseError
		}
		for _, entry := range errResponse.Issue[0].Location {
			if entries := entryRegex.FindStringSubmatch(entry); len(entries) > 1 {
				i, err := strconv.Atoi(entries[1])
				if err != nil {
					return storeResp, err
				}
				invalid[i] = msgs[i]
			}
		}
	}
	if len(invalid) > 0 {
		storeResp.Failed = invalid
		err = ErrBatchErrors
	}
	return storeResp, err
}

func replaceScaryCharacters(msg *Resource) {
	if len(msg.Custom) == 0 {
		return
	}
	stringCustom := strings.Replace(string(msg.Custom), "\\\\", "[bsl]", -1)

	for s, r := range scaryMap {
		stringCustom = strings.Replace(stringCustom, s, r, -1)
	}
	msg.Custom = []byte(stringCustom)
}
