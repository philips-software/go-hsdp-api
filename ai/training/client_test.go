package training_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/ai/training"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

var (
	muxIAM         *http.ServeMux
	serverIAM      *httptest.Server
	muxIDM         *http.ServeMux
	serverIDM      *httptest.Server
	muxTraining    *http.ServeMux
	serverTraining *httptest.Server

	iamClient        *iam.Client
	trainingClient   *training.Client
	trainingTenantID = "48a0183d-a588-41c2-9979-737d15e9e860"
	userUUID         = "e7fecbb2-af8c-47c9-a662-5b046e048bc5"
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)
	muxTraining = http.NewServeMux()
	serverTraining = httptest.NewServer(muxTraining)

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
    "managingOrganization": "`+trainingTenantID+`",
    "organizationList": [
      {
        "organizationId": "`+trainingTenantID+`",
        "permissions": [
          "USER.READ",
          "GROUP.WRITE",
          "DEVICE.READ",
          "CLIENT.SCOPES",
          "AMS_ACCESS.ALL",
          "PKI_CRL_CONFIGURATION.READ",
          "PKI_CERT.ISSUE",
          "PKI_CERT.READ",
          "PKI_CERTS.LIST",
 		  "PKI_CERTROLE.LIST",
   		  "PKI_CERTROLE.READ",
  		  "PKI_URLS.READ",
		  "PKI_CRL.ROTATE",
   		  "PKI_CRL.CONFIGURE",
	      "PKI_CERT.SIGN",
          "PKI_CERT.REVOKE",
          "PKI_URLS.CONFIGURE"
        ],
        "organizationName": "PawneeOrg",
        "groups": [
          "AdminGroup"
        ],
        "roles": [
          "ADMIN",
          "PKIROLE"
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

	trainingClient, err = training.NewClient(iamClient, &ai.Config{
		BaseURL:        serverTraining.URL,
		OrganizationID: trainingTenantID,
	})
	if !assert.Nilf(t, err, "failed to create notificationClient: %v", err) {
		return func() {
		}
	}

	return func() {
		serverIAM.Close()
		serverIDM.Close()
		serverTraining.Close()
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

func TestEndpoint(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	endpoint := serverTraining.URL + "/" + path.Join("analyze", "training", trainingTenantID)

	assert.Equal(t, endpoint, trainingClient.GetEndpointURL())
	err := trainingClient.SetEndpointURL(endpoint)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, endpoint, trainingClient.GetEndpointURL())
}

func TestMethodNotAllowed(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxTraining.HandleFunc("/analyze/training/"+trainingTenantID+"/ComputeEnvironment", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	_, resp, err := trainingClient.ComputeEnvironment.CreateComputeEnvironment(ai.ComputeEnvironment{
		ResourceType: "ComputeEnvironment",
		Name:         "test",
		Image:        "test",
		Description:  "test",
	})
	assert.NotNil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestDebug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	trainingClient, err = training.NewClient(iamClient, &ai.Config{
		BaseURL:        serverTraining.URL,
		DebugLog:       tmpfile.Name(),
		OrganizationID: "xxx",
	})
	if !assert.Nil(t, err) {
		return
	}

	defer trainingClient.Close()
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up

	err = iamClient.Login("username", "password")
	if !assert.Nil(t, err) {
		return
	}
}
