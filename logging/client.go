package logging

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/hsdp/go-signer"
)

var (
	LOG_TIME_FORMAT = "2006-01-02T15:04:05.000Z07:00"
	TIME_FORMAT     = time.RFC3339
	uuidRegex       = regexp.MustCompile(`[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+`)
	versionRegex    = regexp.MustCompile(`^(\d+\.)?(\d+){1}$`)
)

type Config struct {
	SharedKey    string
	SharedSecret string
	BaseURL      string
	ProductKey   string
}

type Client struct {
	config     Config
	url        *url.URL
	httpClient *http.Client
	httpSigner *signer.Signer
	debug      bool
}

func NewClient(httpClient *http.Client, config Config) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
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
		logger.debug = true
	}
	return &logger, nil
}

func (l *Client) Post(msgs []Resource, count int) (err error, sent int, invalid []Resource) {
	var b Bundle

	b.ResourceType = "Bundle"
	b.ProductKey = l.config.ProductKey
	b.Entry = make([]Element, count, count)
	b.Type = "transaction"

	j := 0
	for i := 0; i < count; i++ {
		msg := msgs[i]
		if !msg.Valid() {
			if invalid == nil {
				invalid = make([]Resource, count, count)
			}
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
	b.Total = j

	req := &http.Request{
		Method:     http.MethodPost,
		URL:        l.url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       l.url.Host,
	}

	bodyBytes, err := json.Marshal(b)
	if err != nil {
		return err, 0, invalid
	}
	bodyReader := bytes.NewReader(bodyBytes)
	req.Body = ioutil.NopCloser(bodyReader)
	req.ContentLength = int64(bodyReader.Len())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Version", "1")
	l.httpSigner.SignRequest(req)

	if l.debug {
		dumped, _ := httputil.DumpRequest(req, true)
		fmt.Printf("REQUEST: %s\n", string(dumped))
	}
	resp, err := l.httpClient.Do(req)

	if l.debug {
		if resp != nil {
			fmt.Fprintf(os.Stderr, "Response status: HTTP %d\n", resp.StatusCode)
		} else {
			fmt.Fprintf(os.Stderr, "Error sending response: %s\n", err)
		}
	}
	if err != nil {
		return err, 0, invalid
	}
	return nil, b.Total, invalid
}
