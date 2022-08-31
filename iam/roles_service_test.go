package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleName := "TESTROLE"
	roleDescription := "Role description"
	managingOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	roleID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
				"name": "`+roleName+`",
				"description": "`+roleDescription+`",
				"managingOrganization": "`+managingOrgID+`",
				"id": "`+roleID+`"
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Role/"+roleID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
						"name": "`+roleName+`",
						"description": "`+roleDescription+`",
						"managingOrganization": "`+managingOrgID+`",
						"id": "`+roleID+`"
			}`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var r Role
	r.Name = roleName
	r.Description = roleDescription
	r.ManagingOrganization = managingOrgID

	role, resp, err := client.Roles.CreateRole(roleName, roleDescription, managingOrgID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != http.StatusCreated {
		t.Errorf("Expected HTTP created. Got: %d", resp.StatusCode())
	}
	if role.Name != roleName {
		t.Errorf("Expected role name: %s, Got: %s", roleName, role.Name)
	}
	if role.ManagingOrganization != managingOrgID {
		t.Errorf("Expected role managing Org ID: %s, Got: %s", managingOrgID, role.ManagingOrganization)
	}

	role, resp, err = client.Roles.GetRoleByID(roleID)
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode())
	}
	if role == nil {
		t.Errorf("Expected role to be found, Got: %v", err)
		return
	}
	if role.ID != roleID {
		t.Errorf("Expected to find role with ID: %s, Got: %s", roleID, role.ID)
	}
	_, resp, err = client.Roles.DeleteRole(*role)
	if !assert.Nil(t, err) {
		return
	}
	if resp.StatusCode() != http.StatusNoContent {
		t.Errorf("Expected HTTP no content Got: %d", resp.StatusCode())
	}

}

func roleActionSuccessHandler(message string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
                      "resourceType": "OperationOutcome",
                      "issue": [
                        {
                          "severity": "information",
                          "code": "informational",
                          "details": {
                            "coding": {},
                            "text": "`+message+`"
                          },
                          "diagnostics": "`+message+`"
                        }
                      ]
                    }`)
	}
}

func TestRolePermissionActions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleID := "678abffc-dea1-11e8-9e14-6a0002b8cb70"
	assignMessage := "Permission(s) assigned successfully"
	removeMessage := "Permission(s) removed successfully"

	muxIDM.HandleFunc("/authorize/identity/Role/"+roleID+"/$assign-permission", roleActionSuccessHandler(assignMessage))
	muxIDM.HandleFunc("/authorize/identity/Role/"+roleID+"/$remove-permission", roleActionSuccessHandler(removeMessage))
	var role Role
	role.ID = roleID

	_, resp, err := client.Roles.AddRolePermission(role, "GROUP.READ")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	_, resp, err = client.Roles.RemoveRolePermission(role, "GROUP.READ")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetRolesByGroupID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleName := "TESTROLE"
	roleDescription := "Role description"
	managingOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	roleID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	groupID := "3c7a0274-169e-4ea9-ad91-252cc4022605"
	muxIDM.HandleFunc("/authorize/identity/Role", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("groupId") != groupID {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
		    "total": 1,
		    "entry": [{
				"name": "`+roleName+`",
				"description": "`+roleDescription+`",
				"managingOrganization": "`+managingOrgID+`",
				"id": "`+roleID+`"
			}]
		    }`)
	})
	roles, resp, err := client.Roles.GetRolesByGroupID(groupID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode())
	}
	if len(*roles) != 1 {
		t.Errorf("Expected 1 role")
	}
}

func TestGetRolePermissions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	permissionName := "TESTPERMISSION"
	roleDescription := "Role description"
	roleID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Permission", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, roleID, r.URL.Query().Get("roleId")) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
		    "total": 1,
		    "entry": [{
				"name": "`+permissionName+`",
				"description": "`+roleDescription+`",
				"category": "FOO",
				"type": "GLOBAL",
				"id": "`+roleID+`"
			}]
		    }`)
	})
	var role Role
	role.ID = roleID
	permissions, resp, err := client.Roles.GetRolePermissions(role)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Nil(t, err)
	if !assert.NotNil(t, permissions) {
		return
	}
	assert.Contains(t, *permissions, permissionName)
}
