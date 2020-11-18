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
