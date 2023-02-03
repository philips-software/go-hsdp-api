package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 207, 304:
		return nil
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		data = []byte(err.Error())
	}
	if data == nil {
		data = []byte("empty")
	}
	r.Body = io.NopCloser(bytes.NewBuffer(data)) // Preserve body
	requestURI := ""
	if r.Request.URL != nil {
		requestURI = r.Request.URL.RequestURI()
	}
	return fmt.Errorf("%s %s: StatusCode %d, Body: %s", r.Request.Method, requestURI, r.StatusCode, string(data))
}
