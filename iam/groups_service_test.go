package iam

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	_, _, err := client.Groups.CreateGroup(g)
	assert.NotNil(t, err)

	g.Name = groupName
	group, resp, err := client.Groups.CreateGroup(g)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, groupName, group.Name)
	assert.Equal(t, managingOrgID, group.ManagingOrganization)

	group.Description = updateDescription
	group, resp, err = client.Groups.UpdateGroup(*group)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	group, resp, err = client.Groups.GetGroupByID(group.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, groupID, group.ID)

	ok, resp, err := client.Groups.DeleteGroup(*group)
	assert.True(t, ok)
	assert.Nil(t, err)
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
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
	assert.True(t, ok)
	assert.Nil(t, err)
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
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
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, resp, err = client.Groups.AddMembers(group, "foo")
	assert.NotNil(t, resp)
	assert.False(t, ok)
	assert.NotNil(t, err)
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
	assert.NotNil(t, resp)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Groups.RemoveMembers(group, "foo")
	assert.NotNil(t, resp)
	assert.False(t, ok)
	assert.NotNil(t, err)
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
	assert.NotNil(t, resp)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
	assert.NotNil(t, resp)
	assert.NotNil(t, roles)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(*roles))
	assert.Equal(t, roleID, (*roles)[0].ID)
	assert.Equal(t, "TDRALL", (*roles)[0].Name)
}

func TestAddServices(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceID := "f5fe538f-c3b5-4454-8774-cd3789f59b9a"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$assign", func(w http.ResponseWriter, r *http.Request) {
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
				MemberType string   `json:"memberType"`
				Value      []string `json:"value"`
			}
			err = json.Unmarshal(body, &addRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if addRequest.MemberType != "SERVICE" {
				t.Errorf("Expected SERVICE MemberType, got: %s", addRequest.MemberType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Value); l != 1 {
				t.Errorf("Expected 1 value, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if n := addRequest.Value[0]; n != serviceID {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"services/`+serviceID+`"
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
	ok, resp, err := client.Groups.AddServices(group, serviceID)
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, resp, err = client.Groups.AddServices(group, "foo")
	assert.NotNil(t, resp)
	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRemoveServices(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceID := "f5fe538f-c3b5-4454-8774-cd3789f59b9a"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$remove", func(w http.ResponseWriter, r *http.Request) {
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
				MemberType string   `json:"memberType"`
				Value      []string `json:"value"`
			}
			err = json.Unmarshal(body, &addRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if addRequest.MemberType != "SERVICE" {
				t.Errorf("Expected SERVICE MemberType, got: %s", addRequest.MemberType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if l := len(addRequest.Value); l != 1 {
				t.Errorf("Expected 1 value, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r := addRequest.Value[0]; r != serviceID {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"services/`+r+`"
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
	ok, resp, err := client.Groups.RemoveServices(group, serviceID)
	assert.NotNil(t, resp)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err = client.Groups.RemoveServices(group, "foo")
	assert.NotNil(t, resp)
	assert.False(t, ok)
	assert.NotNil(t, err)
}
