package audit_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/audit"

	"github.com/google/fhir/go/jsonformat"

	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

var (
	muxAudit    *http.ServeMux
	serverAudit *httptest.Server

	iamClient    *iam.Client
	auditClient  *audit.Client
	cdrOrgID     = "48a0183d-a588-41c2-9979-737d15e9e860"
	userUUID     = "e7fecbb2-af8c-47c9-a662-5b046e048bc5"
	timeZone     = "UTC"
	token        string
	refreshToken string
	ma           *jsonformat.Marshaller
	um           *jsonformat.Unmarshaller
)

func setup(t *testing.T) func() {
	var err error
	muxAudit = http.NewServeMux()
	serverAudit = httptest.NewServer(muxAudit)

	auditClient, err = audit.NewClient(nil, &audit.Config{
		AuditBaseURL: serverAudit.URL,
		SharedKey:    "foo",
		SharedSecret: "bar",
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid client")
	}
	ma, err = jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
	if !assert.Nil(t, err) {
		t.Fatalf("failed to create marshaller")
	}
	um, err = jsonformat.NewUnmarshaller("Europe/Amsterdam", jsonformat.STU3)
	if !assert.Nil(t, err) {
		t.Fatalf("failed to create unmarshaller")
	}

	return func() {
		serverAudit.Close()
	}
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	auditClient.CreateAuditEvent(nil)

	defer auditClient.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	// TODO: trigger call here

	fi, err := tmpfile.Stat()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, fi.Size(), "Expected something to be written to DebugLog")
}
