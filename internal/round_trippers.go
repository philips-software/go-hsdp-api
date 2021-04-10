package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

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
}

func NewLoggingRoundTripper(next http.RoundTripper, logFile *os.File) *LoggingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &LoggingRoundTripper{
		next:    next,
		logFile: logFile,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	localID := rt.id
	rt.id++

	if rt.logFile != nil {
		out := ""
		dumped, err := httputil.DumpRequest(req, true)
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %d] --- Request dump error: %v\n", localID, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %d] --- Request start ---\n%s\n[go-hsdp-api %d] Request end ---\n", localID, string(dumped), localID)
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
		if err != nil {
			out = fmt.Sprintf("[go-hsdp-api %d] --- Response dump error: %v\n", localID, err)
		} else {
			out = fmt.Sprintf("[go-hsdp-api %d] --- Response start ---\n%s\n[go-hsdp-api %d] --- Response end ---\n", localID, string(dumped), localID)
		}
		_, _ = rt.logFile.WriteString(out)
	}

	return resp, err
}
