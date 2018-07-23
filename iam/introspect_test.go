package iam

import (
	"io"
	"net/http"
	"testing"
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
		io.WriteString(w, `{
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
						]
					},
					{
						"organizationId": "e5550a19-b6d9-4a9b-ac3c-10ba817776d4",
						"permissions": [
							"ROLE.WRITE",
							"PERMISSION.READ",
							"CLIENT.READ",
							"CLIENT.WRITE"
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
						]
					}
				]
			},
			"client_id": "SomeClient",
			"token_type": "Bearer",
			"identity_type": "user"
		}`)
	})

	_, resp, err := client.Introspect()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if !client.HasPermissions(orgID, "SERVICE.SCOPE", "ORGANIZATION.MFA") {
		t.Errorf("Expected SERVICE.SCOPE and ORGANIZATION.MFA to be there")
	}
	if client.HasPermissions("bogus", "SERVICE.SCOPE") {
		t.Errorf("Bogus orgID should not return any permissions")
	}
}
