package pki_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/has"

	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

var (
	muxIAM    *http.ServeMux
	serverIAM *httptest.Server
	muxIDM    *http.ServeMux
	serverIDM *httptest.Server
	muxHAS    *http.ServeMux
	serverHAS *httptest.Server

	iamClient *iam.Client
	hasClient *has.Client
	hasOrgID  = "48a0183d-a588-41c2-9979-737d15e9e860"
	userUUID  = "e7fecbb2-af8c-47c9-a662-5b046e048bc5"
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)
	muxHAS = http.NewServeMux()
	serverHAS = httptest.NewServer(muxHAS)

	var err error
	iamClient, err = iam.NewClient(nil, &iam.Config{
		OAuth2ClientID: "TestClient",
		OAuth2Secret:   "Secret",
		SharedKey:      "SharedKey",
		SecretKey:      "SecretKey",
		IAMURL:         serverIAM.URL,
		IDMURL:         serverIDM.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create iamCleitn: %v", err)
	}
	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	muxIAM.HandleFunc("/authorize/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
    "scope": "mail tdr.contract tdr.dataitem",
    "access_token": "`+token+`",
    "refresh_token": "31f1a449-ef8e-4bfc-a227-4f2353fde547",
    "expires_in": 1799,
    "token_type": "Bearer"
}`)
	})
	muxIAM.HandleFunc("/authorize/oauth2/introspect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "active": true,
  "scope": "auth_iam_organization auth_iam_introspect mail openid profile cn",
  "username": "ronswanson",
  "exp": 1592073485,
  "sub": "`+userUUID+`",
  "iss": "https://iam-client-test.us-east.philips-healthsuite.com/oauth2/access_token",
  "organizations": {
    "managingOrganization": "`+hasOrgID+`",
    "organizationList": [
      {
        "organizationId": "`+hasOrgID+`",
        "permissions": [
          "ANALYZE_WORKFLOW.ALL",
          "ANALYZE_DATAPROC.ALL",
          "USER.READ",
          "GROUP.WRITE",
          "DEVICE.READ",
          "CLIENT.SCOPES",
          "AMS_ACCESS.ALL",
          "HAS_RESOURCE.ALL",
          "HAS_SESSION.ALL"
        ],
        "organizationName": "PawneeOrg",
        "groups": [
          "AdminGroup"
        ],
        "roles": [
          "ADMIN",
          "HASROLE"
        ]
      }
    ]
  },
  "client_id": "testclientid",
  "token_type": "Bearer",
  "identity_type": "user"
}`)
	})

	// Login immediately so we can create tdrClient
	err = iamClient.Login("username", "password")
	assert.Nil(t, err)

	hasClient, err = has.NewClient(iamClient, &has.Config{
		HASURL: serverHAS.URL,
		OrgID:  hasOrgID,
	})
	assert.Nilf(t, err, "failed to create hasClient: %v", err)

	return func() {
		serverIAM.Close()
		serverIDM.Close()
		serverHAS.Close()
	}
}

func TestLogin(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	token := "44d20214-7879-4e35-923d-f9d4e01c9746"

	err := iamClient.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token, iamClient.Token())
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	orgID := "48a0183d-a588-41c2-9979-737d15e9e860"

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	hasClient, err = has.NewClient(iamClient, &has.Config{
		HASURL:   serverHAS.URL,
		Debug:    true,
		OrgID:    orgID,
		DebugLog: tmpfile.Name(),
	})
	if !assert.Nil(t, err) {
		return
	}

	defer hasClient.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	err = iamClient.Login("username", "password")
	if !assert.Nil(t, err) {
		return
	}

	_, _, _ = hasClient.Resources.GetResources(&has.ResourceOptions{})
	_, _, _ = hasClient.Resources.CreateResource(has.Resource{})
	_, _, _ = hasClient.Resources.DeleteResources(&has.ResourceOptions{})

	fi, err := tmpfile.Stat()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, fi.Size(), "Expected something to be written to DebugLog")
}
