package signer

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC)
}

func TestSigner(t *testing.T) {
	signer, _ := New("foo", "bar", "", fixedTime)
	req, _ := http.NewRequest("GET", "https://example.com/path", nil)

	signer.SignRequest(req)

	signedDate := req.Header.Get(SIGNED_DATE_HEADER)
	signature := req.Header.Get(AUTHORIZATION_HEADER)

	nowFormatted := fixedTime().UTC().Format(TIME_FORMAT)

	if signedDate != nowFormatted {
		t.Error(fmt.Sprintf("Signature mismatch: %s != %s", signedDate, nowFormatted))
	}
	if signature != "HmacSHA256;Credential:foo;SignedHeaders:SignedDate;Signature:mws6Zf5yd8e2dhiCR0fMVyaisvLliNNqnCWpyy1am08=" {
		t.Error(fmt.Sprintf("Invalid signture: %s", signature))
	}
}

func TestValidator(t *testing.T) {
	signer, _ := New("foo", "bar", "", fixedTime)
	req, _ := http.NewRequest("GET", "https://example.com/path", nil)

	signer.SignRequest(req)

	valid, err := signer.ValidateRequest(req)
	if !valid {
		t.Error(fmt.Sprintf("Validation failed: %s", err))
	}
}
