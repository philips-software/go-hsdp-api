package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type filter struct {
	Regex   *regexp.Regexp
	Replace string
}

const (
	Amos = "SOMETIMES_YOU_GOT_TO_STOP_THINKING_ABOUT_SOMETHING_TO_FIGURE_IT_OUT"
)

var filterList = []filter{
	{regexp.MustCompile(`Authorization: (.*)\n`), "Authorization: [sensitive]\n"},
	{regexp.MustCompile(`X-User-Access-Token: (.*)\n`), "X-User-Access-Token: [sensitive]\n"},
	{regexp.MustCompile(`password=[\w%]+`), "password=sensitive"},
	{regexp.MustCompile(`"refresh_token":"[^"]+"`), `"refresh_token":"[sensitive]"`},
	{regexp.MustCompile(`"access_token":"[^"]+"`), `"access_token":"[sensitive]"`},
	{regexp.MustCompile(`"id_token":"[^"]+"`), `"id_token":"[sensitive]"`},
	{regexp.MustCompile(`"token":"[^"]+"`), `"token":"[sensitive]"`},
	{regexp.MustCompile(`token=\w+`), `token=sensitive`},
	{regexp.MustCompile(`id_token_hint=\w+`), `id_token_hint=sensitive`},
	{regexp.MustCompile(`assertion=[\w%.-]+`), `assertion=sensitive`},
	{regexp.MustCompile(`"privateKey":\s*"[^"]+"`), `"privateKey": "[sensitive]"`},
	{regexp.MustCompile(`"productKey":\s*"[^"]+"`), `"productKey": "[sensitive]"`},
	{regexp.MustCompile(`"auth":\s*"[^"]+"`), `"auth": "[sensitive]"`},
}

type HeaderFunc func(req *http.Request) error

type HeaderRoundTripper struct {
	next            http.RoundTripper
	Header          http.Header
	HeaderFunctions []HeaderFunc
}

func NewHeaderRoundTripper(next http.RoundTripper, Header http.Header, functions ...HeaderFunc) *HeaderRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &HeaderRoundTripper{
		next:            next,
		Header:          Header,
		HeaderFunctions: functions,
	}
}

func (rt *HeaderRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.Header != nil {
		for k, v := range rt.Header {
			req.Header[k] = v
		}
	}
	for _, f := range rt.HeaderFunctions {
		_ = f(req)
	}
	return rt.next.RoundTrip(req)
}

type LoggingRoundTripper struct {
	next    http.RoundTripper
	logFile *os.File
	id      int64
	prefix  string
	debug   bool
}

func NewLoggingRoundTripper(next http.RoundTripper, logFile *os.File) *LoggingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &LoggingRoundTripper{
		next:    next,
		logFile: logFile,
		prefix:  uuid.New().String(),
		debug:   os.Getenv(Amos) == "true",
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// We should use a mutex here
	localID := rt.id
	rt.id++

	id := fmt.Sprintf("%s-%05d", rt.prefix, localID)
	if rt.logFile != nil {
		now := time.Now().UTC().Format(time.RFC3339Nano)
		out := ""
		dumped, err := httputil.DumpRequest(req, true)
		filtered := string(dumped)
		if !rt.debug {
			for _, f := range filterList {
				filtered = f.Regex.ReplaceAllString(filtered, f.Replace)
			}
		}
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %s %s] --- request dump error: %v\n", id, now, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %s %s] --- request start ---\n%s\n[go-hsdp-api %s %s] request end ---\n", id, now, filtered, id, now)
		}
		_, _ = rt.logFile.WriteString(out)
	}

	resp, err = rt.next.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if rt.logFile != nil {
		now := time.Now().UTC().Format(time.RFC3339Nano)
		out := ""
		dumped, err := httputil.DumpResponse(resp, true)
		filtered := string(dumped)
		if !rt.debug {
			for _, f := range filterList {
				filtered = f.Regex.ReplaceAllString(filtered, f.Replace)
			}
		}
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %s %s] --- response dump error: %v\n", id, now, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %s %s] --- response start ---\n%s\n[go-hsdp-api %s %s] --- response end ---\n", id, now, filtered, id, now)
		}
		_, _ = rt.logFile.WriteString(out)
	}

	return resp, err
}
