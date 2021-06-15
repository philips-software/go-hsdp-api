package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"

	"github.com/google/uuid"
)

type filter struct {
	Regex   *regexp.Regexp
	Replace string
}

var filterList = []filter{
	{regexp.MustCompile(`Authorization: (.*)\n`), "Authorization: [sensitive]\n"},
	{regexp.MustCompile(`password=\w+`), "password=sensitive"},
	{regexp.MustCompile(`"refresh_token":"[^"]+"`), `"refresh_token":"[sensitive]"`},
	{regexp.MustCompile(`"access_token":"[^"]+"`), `"access_token":"[sensitive]"`},
	{regexp.MustCompile(`"id_token":"[^"]+"`), `"id_token":"[sensitive]"`},
	{regexp.MustCompile(`"token":"[^"]+"`), `"token":"[sensitive]"`},
	{regexp.MustCompile(`token=\w+`), `token=sensitive`},
}

type HeaderRoundTripper struct {
	next   http.RoundTripper
	Header http.Header
}

func NewHeaderRoundTripper(next http.RoundTripper, Header http.Header) *HeaderRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &HeaderRoundTripper{
		next:   next,
		Header: Header,
	}
}

func (rt *HeaderRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.Header != nil {
		for k, v := range rt.Header {
			req.Header[k] = v
		}
	}
	return rt.next.RoundTrip(req)
}

type LoggingRoundTripper struct {
	next    http.RoundTripper
	logFile *os.File
	id      int64
	prefix  string
}

func NewLoggingRoundTripper(next http.RoundTripper, logFile *os.File) *LoggingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &LoggingRoundTripper{
		next:    next,
		logFile: logFile,
		prefix:  uuid.New().String(),
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	localID := rt.id
	rt.id++

	id := fmt.Sprintf("%s-%d", rt.prefix, localID)
	if rt.logFile != nil {
		out := ""
		dumped, err := httputil.DumpRequest(req, true)
		filtered := string(dumped)
		for _, f := range filterList {
			filtered = f.Regex.ReplaceAllString(filtered, f.Replace)
		}
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %s] --- Request dump error: %v\n", id, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %s] --- Request start ---\n%s\n[go-hsdp-api %s] Request end ---\n", id, filtered, id)
		}
		_, _ = rt.logFile.WriteString(out)
	}

	resp, err = rt.next.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if rt.logFile != nil {
		out := ""
		dumped, err := httputil.DumpResponse(resp, true)
		filtered := string(dumped)
		for _, f := range filterList {
			filtered = f.Regex.ReplaceAllString(filtered, f.Replace)
		}
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %s] --- Response dump error: %v\n", id, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %s] --- Response start ---\n%s\n[go-hsdp-api %s] --- Response end ---\n", id, filtered, id)
		}
		_, _ = rt.logFile.WriteString(out)
	}

	return resp, err
}
