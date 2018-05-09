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

	"github.com/m4rw3r/uuid"
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

type Logger struct {
	config     Config
	url        *url.URL
	httpClient *http.Client
	debug      bool
}

func NewClient(httpClient *http.Client, config Config) (*Logger, error) {
	var logger Logger

	logger.config = config
	logger.httpClient = httpClient
	url, err := url.Parse(config.BaseURL + "/core/log/LogEvent")
	if err != nil {
		return nil, err
	}
	logger.url = url
	if os.Getenv("DEBUG") == "true" {
		logger.debug = true
	}
	return &logger, nil
}

func (l *Logger) Post(msgs []DHPLogMessage, count int) error {
	var b Bundle

	// Bundle
	b.ResourceType = "Bundle"
	b.ProductKey = l.config.ProductKey
	b.Type = "transaction"
	b.Total = count
	b.Entry = make([]Element, count, count)

	for i := 0; i < count; i++ {
		// Element
		var e Element
		msg := msgs[i]
		e.Resource.ApplicationInstance = msg.ApplicationInstance
		e.Resource.ApplicationName = msg.ApplicationName
		e.Resource.Category = msg.Category
		e.Resource.ApplicationVersion = msg.ApplicationVersion
		if e.Resource.ApplicationVersion == "" {
			e.Resource.ApplicationVersion = "0.0.0"
		}
		e.Resource.Component = "PHS"
		if msg.Component != "" && !(msg.Component == "DHP" || msg.Component == "CPH") {
			e.Resource.Component = msg.Component
		}
		if msg.EventID == "" {
			msg.EventID = "1"
		}
		e.Resource.EventID = msg.EventID
		e.Resource.LogTime = msg.LogTime
		id, _ := uuid.V4()
		e.Resource.ID = id.String()
		e.Resource.OriginatingUser = msg.OriginatingUser
		e.Resource.ServerName = msg.ServerName
		if e.Resource.ServerName == "" {
			e.Resource.ServerName = "not-set"
		}
		e.Resource.ServiceName = msg.ServiceName
		e.Resource.Severity = msg.Severity
		e.Resource.OriginatingUser = msg.OriginatingUser
		if e.Resource.OriginatingUser == "" {
			e.Resource.OriginatingUser = "not-specified"
		}
		if uuidRegex.MatchString(msg.TransactionID) {
			e.Resource.TransactionID = msg.TransactionID
		} else {
			trns, _ := uuid.V4()
			e.Resource.TransactionID = trns.String()
		}
		e.Resource.LogData.Message = base64.StdEncoding.EncodeToString([]byte(msg.LogData.Message))
		e.Resource.ResourceType = "LogEvent"

		b.Entry[i] = e
	}

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
		return err
	}
	bodyReader := bytes.NewReader(bodyBytes)
	req.Body = ioutil.NopCloser(bodyReader)
	req.ContentLength = int64(bodyReader.Len())
	req.Header.Set("Content-Type", "application/json")

	if l.debug {
		dumped, _ := httputil.DumpRequest(req, true)
		fmt.Printf("REQUEST: %v\n", dumped)
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
		return err
	}
	return nil
}
