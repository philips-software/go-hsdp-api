package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleID := "168cdeae-2539-45b0-b18c-89ae32f1ea15"

	muxIDM.HandleFunc("/authorize/identity/Permission", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"total": 3,
			"entry": [
				{
					"name": "SERVICE.SCOPE",
					"category": "IAM",
					"type": "GLOBAL",
					"id": "f1c8b67a-e652-4a91-abb1-0b5d032948dd"
				},
				{
					"name": "ROLE.WRITE",
					"category": "IAM",
					"type": "GLOBAL",
					"id": "11615a64-34dd-4ada-be73-b30a0acb8769"
				},
				{
					"name": "ORGANIZATION.MFA",
					"category": "IAM",
					"type": "GLOBAL",
					"id": "363f6953-158c-4122-af76-b997f259c4af"
				}
			]
		}`)
	})

	permissions, resp, err := client.Permissions.GetPermissions(GetPermissionOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if len(*permissions) != 3 {
		t.Errorf("Expected 3 permissions, Got: %d", len(*permissions))
	}
	permissions, resp, err = client.Permissions.GetPermissionsByRoleID(roleID)
	if len(*permissions) != 3 {
		t.Errorf("Expected 3 permissions, Got: %d", len(*permissions))
	}
}

func TestGetPermissionByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	uuid := "f1c8b67a-e652-4a91-abb1-0b5d032948dd"

	muxIDM.HandleFunc("/authorize/identity/Permission", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
			"total": 3,
			"entry": [
				{
					"name": "SOME.SCOPE",
					"category": "IAM",
					"type": "GLOBAL",
					"id": "`+uuid+`"
				}
			]
		}`)
	})

	permission, resp, err := client.Permissions.GetPermissionByID(uuid)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success")
	}
	if permission.ID != uuid {
		t.Errorf("Expected Permission with ID: %s, Got: %s", uuid, permission.ID)
	}
}
