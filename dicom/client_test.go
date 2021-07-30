package dicom_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/dicom"

	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

var (
	muxIAM      *http.ServeMux
	serverIAM   *httptest.Server
	muxIDM      *http.ServeMux
	serverIDM   *httptest.Server
	muxDICOM    *http.ServeMux
	serverDICOM *httptest.Server

	iamClient   *iam.Client
	dicomClient *dicom.Client
	cdrOrgID    = "48a0183d-a588-41c2-9979-737d15e9e860"
	userUUID    = "e7fecbb2-af8c-47c9-a662-5b046e048bc5"
)

func setup(t *testing.T) func() {
	muxIAM = http.NewServeMux()
	serverIAM = httptest.NewServer(muxIAM)
	muxIDM = http.NewServeMux()
	serverIDM = httptest.NewServer(muxIDM)
	muxDICOM = http.NewServeMux()
	serverDICOM = httptest.NewServer(muxDICOM)

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
    "managingOrganization": "`+cdrOrgID+`",
    "organizationList": [
      {
        "organizationId": "`+cdrOrgID+`",
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

	// Login immediately so we can create dicomClient
	err = iamClient.Login("username", "password")
	assert.Nil(t, err)

	dicomClient, err = dicom.NewClient(iamClient, &dicom.Config{
		DICOMConfigURL: serverDICOM.URL,
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid client")
	}

	return func() {
		serverIAM.Close()
		serverIDM.Close()
		serverDICOM.Close()
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

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	dicomClient, err = dicom.NewClient(iamClient, &dicom.Config{
		DICOMConfigURL: serverDICOM.URL,
		DebugLog:       tmpfile.Name(),
	})
	if !assert.Nil(t, err) {
		return
	}

	defer dicomClient.Close()
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up

	err = iamClient.Login("username", "password")
	if !assert.Nil(t, err) {
		return
	}

	fi, err := tmpfile.Stat()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, fi.Size(), "Expected something to be written to DebugLog")
}

func TestEndpoints(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	dicomClient, err := dicom.NewClient(iamClient, &dicom.Config{
		DICOMConfigURL: serverDICOM.URL,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, dicomClient) {
		return
	}
	assert.Equal(t, serverDICOM.URL+"/store/dicom/", dicomClient.GetDICOMStoreURL())

}

func TestGeneratedURLs(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	dicomClient, err := dicom.NewClient(iamClient, &dicom.Config{
		DICOMConfigURL: "https: //dss-config-share-tst.foo.io",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, dicomClient) {
		return
	}
	assert.Equal(t, "https: //dss-qido-share-tst.foo.io", dicomClient.GetQIDOURL())
	assert.Equal(t, "https: //dss-stow-share-tst.foo.io", dicomClient.GetSTOWURL())
	assert.Equal(t, "https: //dss-wado-share-tst.foo.io", dicomClient.GetWADOURL())
}

func TestErrorResponse(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxDICOM.HandleFunc("/store/dicom/config/dicom/production/remoteNodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		w.WriteHeader(http.StatusConflict)
		_, _ = io.WriteString(w, `{"error":"something unexpected happened"}`)
	})

	_, resp, err := dicomClient.Config.CreateRemoteNode(dicom.RemoteNode{
		Title: "Some Title here",
		NetworkConnection: dicom.NetworkConnection{
			Port:     31337,
			HostName: "foo.com",
		},
		AETitle: "AE Title here",
	}, nil)
	if !assert.NotNil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, err.Error(), `POST : StatusCode 409, Body: {"error":"something unexpected happened"}`)
}
