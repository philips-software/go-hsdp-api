package iam

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
	g.Description = groupDescription
	g.ManagingOrganization = managingOrgID

	// Test validation
	group, resp, err := client.Groups.CreateGroup(g)
	if err == nil {
		t.Error("Expected validation failure")
	}

	g.Name = groupName
	group, resp, err = client.Groups.CreateGroup(g)
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
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Unexpected EOF from reading request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var assignRequest struct {
				Roles []string `json:"roles"`
			}
			err = json.Unmarshal(body, &assignRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(assignRequest.Roles) != 1 {
				t.Errorf("Expected 1 role, got: %d", len(assignRequest.Roles))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if assignRequest.Roles[0] != roleID {
				t.Errorf("Unexpected role: %s", assignRequest.Roles[0])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
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

func TestAddMembers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$add-members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Unexpected EOF from reading request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var addRequest struct {
				ResourceType string      `json:"resourceType"`
				Parameter    []Parameter `json:"parameter"`
			}
			err = json.Unmarshal(body, &addRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if addRequest.ResourceType != "Parameters" {
				t.Errorf("Expected Parameters resourceType, got: %s", addRequest.ResourceType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Parameter); l != 1 {
				t.Errorf("Expected 1 parameter, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if n := addRequest.Parameter[0].Name; n != "UserIDCollection" {
				t.Errorf("Unexpected parameters: %s", n)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Parameter[0].References); l == 0 {
				t.Errorf("Expected at least 1 reference, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r := addRequest.Parameter[0].References[0].Reference; r != userID {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"users/`+r+`"
							]
						}
					]
				}
				`)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	var group Group
	group.ID = groupID
	ok, resp, err := client.Groups.AddMembers(group, userID)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if !ok {
		t.Errorf("Expected AddMembers to succeed")
	}
	if err != nil {
		t.Errorf("Did not expect error, Got: %v", err)
	}
	ok, resp, err = client.Groups.AddMembers(group, "foo")
	if ok {
		t.Errorf("Expected AddMembers to fail")
	}
	if err == nil {
		t.Errorf("Expected error from AddMembers")
	}
}

func TestRemoveMembers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$remove-members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Unexpected EOF from reading request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var addRequest struct {
				ResourceType string      `json:"resourceType"`
				Parameter    []Parameter `json:"parameter"`
			}
			err = json.Unmarshal(body, &addRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if addRequest.ResourceType != "Parameters" {
				t.Errorf("Expected Parameters resourceType, got: %s", addRequest.ResourceType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Parameter); l != 1 {
				t.Errorf("Expected 1 parameter, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if n := addRequest.Parameter[0].Name; n != "UserIDCollection" {
				t.Errorf("Unexpected parameters: %s", n)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Parameter[0].References); l == 0 {
				t.Errorf("Expected at least 1 reference, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r := addRequest.Parameter[0].References[0].Reference; r != userID {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"users/`+r+`"
							]
						}
					]
				}
				`)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	var group Group
	group.ID = groupID
	ok, resp, err := client.Groups.RemoveMembers(group, userID)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if !ok {
		t.Errorf("Expected RemoveMembers to succeed")
	}
	if err != nil {
		t.Errorf("Did not expect error, Got: %v", err)
	}
	ok, resp, err = client.Groups.AddMembers(group, "foo")
	if ok {
		t.Errorf("Expected RemoveMembers to fail")
	}
	if err == nil {
		t.Errorf("Expected error from RemoveMembers")
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
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Unexpected EOF from reading request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var removeRequest struct {
				Roles []string `json:"roles"`
			}
			err = json.Unmarshal(body, &removeRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(removeRequest.Roles) != 1 {
				t.Errorf("Expected 1 role, got: %d", len(removeRequest.Roles))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if removeRequest.Roles[0] != roleID {
				t.Errorf("Unexpected role: %s", removeRequest.Roles[0])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
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

func TestGetRoles(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	roleID := "c6f4f40d-6585-4dbf-b4bb-cd78bc83e73b"

	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("groupId") != groupID {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
				"total": 1,
				"entry": [
					{
						"name": "TDRALL",
						"managingOrganization": "0d1c477c-46be-4c5c-a53e-51ad86eda38d",
						"id": "`+roleID+`"
					}
				]
			}`)
		}
	})
	var group Group
	var role Role
	group.ID = groupID
	role.ID = roleID
	roles, resp, err := client.Groups.GetRoles(group)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP success Got: %d", resp.StatusCode)
	}
	if err != nil {
		t.Errorf("Did not expect error, Got: %v", err)
	}
	if roles == nil {
		t.Errorf("Expected to find roles")
		return
	}
	if len(*roles) != 1 {
		t.Errorf("Expected to find 1 role, got: %d", len(*roles))
		return
	}
	if (*roles)[0].ID != roleID {
		t.Errorf("Unexpected role ID")
	}
	if (*roles)[0].Name != "TDRALL" {
		t.Errorf("Unexpected role name")
	}
}
