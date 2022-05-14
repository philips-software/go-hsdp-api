package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntrospect(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "b400f634-03ed-4596-bfc1-0b74e5bb1af8"
	orgID := "46323bb4-ebba-4387-a339-252b5aa0755f"

	muxIAM.HandleFunc("/authorize/oauth2/introspect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"active": true,
			"scope": "mail tdr.contract tdr.dataitem",
			"username": "foo.bar@philips.com",
			"exp": 1532350103,
			"sub": "`+userID+`",
			"iss": "hsdp-iam",
			"organizations": {
				"managingOrganization": "`+orgID+`",
				"organizationList": [
					{
						"organizationId": "0aa40c4b-5645-4ff1-b02e-2ec8383ecb29",
						"permissions": [
							"ROLE.WRITE",
							"PERMISSION.READ",
							"LOG.READ",
							"ROLE.READ"
						]
					},
					{
						"organizationId": "c21a69da-c74c-430d-899b-eac0f66e5188",
						"permissions": [
							"WELLCENTIVE_DM.FILEUPLOAD",
							"LOG.READ",
							"ALL.READ",
							"ALL.WRITE"
						]
					},
					{
						"organizationId": "82303602-486d-4372-8002-879637d50807",
						"permissions": [
							"LOG.READ"
						]
					},
					{
						"organizationId": "`+orgID+`",
						"permissions": [
							"APPLICATION.READ",
							"SERVICE.SCOPE",
							"GROUP.WRITE",
							"DEVICE.READ",
							"PROPOSITION.READ",
							"SERVICE.DELETE",
							"GROUP.READ",
							"USER.READ",
							"CLIENT.WRITE",
							"CLIENT.DELETE",
							"ROLE.READ",
							"ROLE.WRITE",
							"PROPOSITION.WRITE",
							"DEVICE.WRITE",
							"PERMISSION.READ",
							"LOG.READ",
							"SERVICE.WRITE",
							"ORGANIZATION.MFA",
							"ORGANIZATION.READ",
							"CLIENT.READ",
							"USER.WRITE",
							"APPLICATION.WRITE",
							"SERVICE.READ",
							"ORGANIZATION.WRITE"
						],
						"effectivePermissions": [
							"APPLICATION.READ",
							"SERVICE.SCOPE",
							"GROUP.WRITE",
							"DEVICE.READ",
							"PROPOSITION.READ",
							"SERVICE.DELETE",
							"GROUP.READ",
							"USER.READ",
							"CLIENT.WRITE",
							"CLIENT.DELETE",
							"ROLE.READ",
							"ROLE.WRITE",
							"PROPOSITION.WRITE",
							"DEVICE.WRITE",
							"PERMISSION.READ",
							"LOG.READ",
							"SERVICE.WRITE",
							"ORGANIZATION.MFA",
							"ORGANIZATION.READ",
							"CLIENT.READ",
							"USER.WRITE",
							"APPLICATION.WRITE",
							"SERVICE.READ",
							"ORGANIZATION.WRITE"
						]
					},
					{
						"organizationId": "e5550a19-b6d9-4a9b-ac3c-10ba817776d4",
						"permissions": [
							"ROLE.WRITE",
							"PERMISSION.READ",
							"CLIENT.READ",
							"CLIENT.WRITE"
						],
						"organizationName": "FirstOrg",
        				"groups": [
          					"S3CredsAdminGroup",
          					"AdminGroup"
        				],
        				"roles": [
          					"ADMIN",
          					"S3CREDSADMINROLE"
        				]
					},
					{
						"organizationId": "f5d34188-57ba-4fe2-afcf-bf8cb57a860b",
						"permissions": [
							"APPLICATION.READ",
							"GROUP.WRITE",
							"PROPOSITION.READ",
							"DEVICE.READ",
							"GROUP.READ",
							"USER.READ",
							"ROLE.READ",
							"ROLE.WRITE",
							"PROPOSITION.WRITE",
							"DEVICE.WRITE",
							"PERMISSION.READ",
							"ORGANIZATION.MFA",
							"ORGANIZATION.READ",
							"USER.WRITE",
							"APPLICATION.WRITE",
							"ORGANIZATION.WRITE"
						],
						"organizationName": "SecondOrg",
        				"groups": [
          					"S3CredsAdminGroup",
          					"AdminGroup"
        				],
        				"roles": [
          					"ADMIN",
          					"S3CREDSADMINROLE"
        				]
					}
				]
			},
			"client_id": "SomeClient",
			"token_type": "Bearer",
			"identity_type": "user"
		}`)
	})

	introspectResponse, resp, err := client.Introspect()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, introspectResponse) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.True(t, client.HasPermissions(orgID, "SERVICE.SCOPE", "ORGANIZATION.MFA"))
	assert.False(t, client.HasPermissions("bogus", "SERVICE.SCOPE"))
	assert.Equal(t, 6, len(introspectResponse.Organizations.OrganizationList))
	assert.Equal(t, "SecondOrg", introspectResponse.Organizations.OrganizationList[5].OrganizationName)
}

func TestIntrospectWithOrgContext(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "b400f634-03ed-4596-bfc1-0b74e5bb1af8"
	orgID := "46323bb4-ebba-4387-a339-252b5aa0755f"

	muxIAM.HandleFunc("/authorize/oauth2/introspect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got ‘%s’", r.Method)
		}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := r.Form.Get("token")
		if token != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		orgContext := r.Form.Get("org_ctx")

		if orgContext != "" && orgContext != "bogus" {
			orgID = orgContext
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"active": true,
			"scope": "mail tdr.contract tdr.dataitem",
			"username": "foo.bar@philips.com",
			"exp": 1532350103,
			"sub": "`+userID+`",
			"iss": "hsdp-iam",
			"organizations": {
				"managingOrganization": "`+orgID+`",
				"organizationList": [
					{
						"organizationId": "`+orgID+`",
						"permissions": [
							"ROLE.WRITE",
							"PERMISSION.READ",
							"LOG.READ",
							"ROLE.READ",
							"SERVICE.SCOPE",
							"ORGANIZATION.MFA"
						],
						"effectivePermissions": [
							"ROLE.WRITE",
							"PERMISSION.READ",
							"LOG.READ",
							"ROLE.READ",
							"SERVICE.SCOPE",
							"ORGANIZATION.MFA"
						]
					}
				]
			},
			"client_id": "SomeClient",
			"token_type": "Bearer",
			"identity_type": "user"
		}`)
	})

	introspectResponse, resp, err := client.Introspect(WithOrgContext(orgID))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, introspectResponse) {
		return
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	assert.True(t, client.HasPermissions(orgID, "SERVICE.SCOPE", "ORGANIZATION.MFA"))
	assert.Equal(t, 1, len(introspectResponse.Organizations.OrganizationList))
	assert.False(t, client.HasPermissions("bogus", "SERVICE.SCOPE"))
}
