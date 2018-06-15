package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestRoleCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleName := "TESTROLE"
	roleDescription := "Role description"
	updateDescription := "Updated description"
	managingOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	roleID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, `{
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
		case "PUT":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
			"name": "`+roleName+`",
			"description": "`+updateDescription+`",
			"managingOrganization": "`+managingOrgID+`",
			"id": "`+roleID+`"
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
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
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP created. Got: %d", resp.StatusCode)
	}
	if role.Name != roleName {
		t.Errorf("Expected role name: %s, Got: %s", roleName, role.Name)
	}
	if role.ManagingOrganization != managingOrgID {
		t.Errorf("Expected role managing Org ID: %s, Got: %s", managingOrgID, role.ManagingOrganization)
	}
	role.Description = updateDescription
	role, resp, err = client.Roles.UpdateRole(role)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	role, resp, err = client.Roles.GetRoleByID(roleID)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if role == nil {
		t.Errorf("Expected role to be found, Got: %v", err)
		return
	}
	if role.ID != roleID {
		t.Errorf("Expected to find role with ID: %s, Got: %s", roleID, role.ID)
	}
	ok, resp, err := client.Roles.DeleteRole(*role)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP no content Got: %d", resp.StatusCode)
	}
	if !ok {
		t.Errorf("Expected role to be deleted")
	}
}
