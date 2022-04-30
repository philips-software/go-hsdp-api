package console

import (
	"fmt"
	"net/http"
	"net/url"
)

func WithHost(host string) OptionFunc {
	return func(req *http.Request) error {
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		q := req.URL.Query()
		q.Set("host", fmt.Sprintf("%s", host))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithQuery(query string) OptionFunc {
	return func(req *http.Request) error {
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		q := req.URL.Query()
		q.Set("query", fmt.Sprintf("%s", query))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithStart(start int64) OptionFunc {
	return func(req *http.Request) error {
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		q := req.URL.Query()
		q.Set("start", fmt.Sprintf("%d", start))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithEnd(end int64) OptionFunc {
	return func(req *http.Request) error {
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		q := req.URL.Query()
		q.Set("end", fmt.Sprintf("%d", end))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithStep(step int64) OptionFunc {
	return func(req *http.Request) error {
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		q := req.URL.Query()
		q.Set("step", fmt.Sprintf("%d", step))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
