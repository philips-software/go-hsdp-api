package pki_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/philips-software/go-hsdp-api/pki"

	"github.com/stretchr/testify/assert"
)

func TestOnboarding(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxPKI.HandleFunc("/core/pki/tenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				return
			}
			var tenant pki.Tenant
			err = json.Unmarshal(body, &tenant)
			if !assert.Nil(t, err) {
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
  "api_endpoint": "`+serverPKI.URL+`/core/pki/api/`+tenant.ServiceParameters.LogicalPath+`"
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	logicalPath := "ron-swanson"
	newTenant := pki.Tenant{
		OrganizationName: "org",
		SpaceName:        "space",
		ServiceName:      "hsdp-pki",
		PlanName:         "standard",
		ServiceParameters: pki.ServiceParameters{
			LogicalPath: logicalPath,
			IAMOrgs:     []string{pkiOrgID},
			CA: pki.CertificateAuthority{
				TTL:        "24h",
				CommonName: "1e100.io",
				KeyType:    "ec",
				KeyBits:    384,
			},
			Roles: []pki.Role{
				{
					Name:            "ec384",
					AllowAnyName:    true,
					AllowIPSans:     true,
					AllowSubdomains: true,
					AllowedURISans: []string{
						"*",
					},
					AllowedOtherSans: []string{
						"*",
					},
					ClientFlag: true,
					Country: []string{
						"NL",
					},
					NotBeforeDuration: "30s",
					EnforceHostnames:  false,
					KeyBits:           384,
					KeyType:           "ec",
					ServerFlag:        true,
					TTL:               "720h",
					UseCSRCommonName:  true,
					UseCSRSans:        true,
				},
			},
		},
	}
	newTenant.ServiceParameters.LogicalPath = logicalPath

	onboarding, resp, err := pkiClient.Tenants.Onboard(newTenant)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, onboarding) {
		return
	}
	assert.True(t, strings.Contains(onboarding.APIEndpoint, logicalPath))
}

func TestOffboarding(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	logicalPath := "ron-swanson"
	muxPKI.HandleFunc("/core/pki/tenant/"+logicalPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	tenant := pki.Tenant{
		ServiceParameters: pki.ServiceParameters{
			LogicalPath: logicalPath,
		},
	}
	ok, resp, err := pkiClient.Tenants.Offboard(tenant)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)

}

func TestRetrieveAndUpdate(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	logicalPath := "ron-swanson"

	muxPKI.HandleFunc("/core/pki/tenant/"+logicalPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				return
			}
			var tenant pki.Tenant
			err = json.Unmarshal(body, &tenant)
			if !assert.Nil(t, err) {
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case "GET":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
  "service_name": "hsdp-pki",
  "plan_name": "standard",
  "organization_name": "org",
  "space_name": "test",
  "service_parameters": {
    "logical_path": "`+logicalPath+`",
    "iam_orgs": [
      "cfb7597f-a812-4fb8-ab32-42b89487fb7e"
    ],
    "ca": {
      "common_name": "andy-test",
      "key_type": "ec",
      "key_bits": 384,
      "ou": "ronswanson",
      "organization": "Pawnee",
      "country": "NL",
      "locality": "Locality",
      "province": "Somsome",
      "ttl": "8640h"
    },
    "roles": [
      {
        "allowed_other_sans": [
          "*"
        ],
        "allowed_uri_sans": [
          "*"
        ],
        "name": "ec384",
        "allow_subdomains": true,
        "allow_any_name": true,
        "enforce_hostnames": false,
        "allow_ip_sans": true,
        "server_flag": true,
        "client_flag": true,
        "ttl": "720h",
        "max_ttl": "8640h",
        "key_type": "ec",
        "key_bits": 384,
        "use_csr_common_name": true,
        "use_csr_sans": true,
        "country": [
          "NL"
        ],
        "not_before_duration": "30s"
      }
    ]
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	tenant, resp, err := pkiClient.Tenants.Retrieve(logicalPath)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, tenant) {
		return
	}
	assert.Equal(t, logicalPath, tenant.ServiceParameters.LogicalPath)

	ok, resp, err := pkiClient.Tenants.Update(*tenant)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.True(t, ok) {
		return
	}
}
