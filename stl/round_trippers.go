package stl

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

type headerRoundTripper struct {
	next   http.RoundTripper
	Header http.Header
}

func newHeaderRoundTripper(next http.RoundTripper, Header http.Header) *headerRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &headerRoundTripper{
		next:   next,
		Header: Header,
	}
}

func (rt *headerRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.Header != nil {
		for k, v := range rt.Header {
			req.Header[k] = v
		}
	}
	return rt.next.RoundTrip(req)
}

type loggingRoundTripper struct {
	next    http.RoundTripper
	logFile *os.File
}

func newLoggingRoundTripper(next http.RoundTripper, logFile *os.File) *loggingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &loggingRoundTripper{
		next:    next,
		logFile: logFile,
	}
}

func (rt *loggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.logFile != nil {
		dumped, _ := httputil.DumpRequest(req, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Request start ---\n%s\n[go-hsdp-api] Request end ---\n", string(dumped))
		_, _ = rt.logFile.WriteString(out)
	}

	resp, err = rt.next.RoundTrip(req)

	if rt.logFile != nil {
		dumped, _ := httputil.DumpResponse(resp, true)
		out := fmt.Sprintf("[go-hsdp-api] --- Response start ---\n%s\n[go-hsdp-api] --- Response end ---\n", string(dumped))
		_, _ = rt.logFile.WriteString(out)
	}

	return resp, err
}
