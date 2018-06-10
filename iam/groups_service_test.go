package iam

import (
	"io"
	"net/http"
	"testing"
)

func TestGroupCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	groupName := "TestGroup"
	groupDescription := "Group description"
	updateDescription := "Updated description"
	managingOrgID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, `{
				"name": "`+groupName+`",
				"description": "`+groupDescription+`",
				"managingOrganization": "`+managingOrgID+`",
				"id": "`+groupID+`"
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
				"resourceType": "bundle",
				"type": "searchset",
				"total": 1,
				"entry": [
					{
						"resource": {
							"resourceType": "Group",
							"groupName": "`+groupName+`",
							"orgId": "`+managingOrgID+`",
							"_id": "`+groupID+`"
						}
					}
				]
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
			"name": "`+groupName+`",
			"description": "`+updateDescription+`",
			"managingOrganization": "`+managingOrgID+`",
			"id": "`+groupID+`"
			}`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var g Group
	g.Name = groupName
	g.Description = groupDescription
	g.ManagingOrganization = managingOrgID

	group, resp, err := client.Groups.CreateGroup(g)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP created. Got: %d", resp.StatusCode)
	}
	if group.Name != groupName {
		t.Errorf("Expected Group name: %s, Got: %s", groupName, group.Name)
	}
	if group.ManagingOrganization != managingOrgID {
		t.Errorf("Expected Group managing Org ID: %s, Got: %s", managingOrgID, group.ManagingOrganization)
	}
	group.Description = updateDescription
	group, resp, err = client.Groups.UpdateGroup(*group)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	group, resp, err = client.Groups.GetGroupByID(group.ID)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if group.ID != groupID {
		t.Errorf("Expected to find group with ID: %s, Got: %s", groupID, group.ID)
	}
	ok, resp, err := client.Groups.DeleteGroup(*group)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP no content Got: %d", resp.StatusCode)
	}
	if !ok {
		t.Errorf("Expected group to be deleted")
	}
}

func TestAssignRole(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$assign-role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
				"resourceType": "OperationOutcome",
				"issue": [
					{
						"severity": "information",
						"code": "informational",
						"details": {
							"text": "Role(s) assigned successfully"
						}
					}
				]
			}`)
		}
	})
	var group Group
	var role Role
	group.ID = groupID
	role.ID = roleID
	ok, resp, err := client.Groups.AssignRole(group, role)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if err != nil {
		t.Errorf("Did not expect error, Got: %v", err)
	}
	if !ok {
		t.Errorf("Expected AssignRole to succeed")
	}
}

func TestRemoveRole(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$remove-role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
				"resourceType": "OperationOutcome",
				"issue": [
					{
						"severity": "information",
						"code": "informational",
						"details": {
							"text": "Role(s) removed successfully"
						}
					}
				]
			}`)
		}
	})
	var group Group
	var role Role
	group.ID = groupID
	role.ID = roleID
	ok, resp, err := client.Groups.RemoveRole(group, role)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if err != nil {
		t.Errorf("Did not expect error, Got: %v", err)
	}
	if !ok {
		t.Errorf("Expected AssignRole to succeed")
	}
}
