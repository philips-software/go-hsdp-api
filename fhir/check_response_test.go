package fhir

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"foo":"bar"}`)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	err := CheckResponse(resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	handlerErr := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"result":"denied"}`)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handlerErr(w, req)

	resp = w.Result()
	resp.Request = req
	err = CheckResponse(resp)
	if err == nil {
		t.Errorf("Expected an error")
		return
	}
	if err.Error() != "GET http://example.com: 401 {result: denied}" {
		t.Errorf("Unexpected error: %v\n", err.Error())
	}

	handlerBroken := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{`)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handlerBroken(w, req)

	resp = w.Result()
	resp.Request = req
	err = CheckResponse(resp)
	if err.Error() != "GET http://example.com: 401 failed to parse unexpected error type: <nil>" {
		t.Errorf("Unexpected error: %v\n", err.Error())
	}
}
