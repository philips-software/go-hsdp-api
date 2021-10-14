package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 207, 304:
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = []byte(err.Error())
	}
	if data == nil {
		data = []byte("empty")
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // Preserve body
	return fmt.Errorf("%s %s: StatusCode %d, Body: %s", r.Request.Method, r.Request.RequestURI, r.StatusCode, string(data))
}
