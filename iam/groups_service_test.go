package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
			_, _ = io.WriteString(w, `{
				"name": "`+groupName+`",
				"description": "`+groupDescription+`",
				"managingOrganization": "`+managingOrgID+`",
				"id": "`+groupID+`"
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"resourceType": "bundle",
				"type": "searchset",
				"total": 1,
				"entry": [
					{
						"resource": {
							"resourceType": "Group",
							"groupDescription": "Group for test",
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
			_, _ = io.WriteString(w, `{
			"name": "`+groupName+`",
			"description": "`+updateDescription+`",
			"managingOrganization": "`+managingOrgID+`",
			"id": "`+groupID+`"
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
			"name": "`+groupName+`",
			"description": "`+groupDescription+`",
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
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	assert.Equal(t, groupName, group.Name)
	assert.Equal(t, managingOrgID, group.ManagingOrganization)

	group.Description = updateDescription
	group, resp, err = client.Groups.UpdateGroup(*group)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	groups, resp, err := client.Groups.GetGroups(&GetGroupOptions{
		OrganizationID: &managingOrgID,
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, groups) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 1, len(*groups))

	group, resp, err = client.Groups.GetGroupByID(group.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, group) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, groupID, group.ID)

	ok, resp, err := client.Groups.DeleteGroup(*group)
	assert.True(t, ok)
	assert.Nil(t, err)
	if ok := assert.NotNil(t, resp); ok {
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	}
}

func TestAssignRole(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	var assignedTotal []string

	roleID := "f5fe538f-c3b5-4454-8774-cd3789f59b9f"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$assign-role", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
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
			assignedTotal = append(assignedTotal, assignRequest.Roles...)

			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
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
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	}
	assert.Len(t, assignedTotal, 1)
}

func TestAddMembers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	var assignedTotal []string

	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$add-members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
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
			for _, a := range addRequest.Parameter[0].References {
				assignedTotal = append(assignedTotal, a.Reference)
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	assert.Len(t, assignedTotal, 0)

	var group Group
	group.ID = groupID
	var users []string
	for i := 0; i < 28; i++ {
		users = append(users, fmt.Sprintf("%s%02d", "f5fe538f-c3b5-4454-8774-cd3789f59b", i))
	}
	ok, resp, err := client.Groups.AddMembers(context.Background(), group, users...)
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	}
	if !assert.Nil(t, err) {
		return
	}
	assert.Len(t, assignedTotal, 28)
	assert.Nil(t, ok)
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
			body, err := io.ReadAll(r.Body)
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
			l := len(addRequest.Parameter[0].References)
			if l == 0 {
				t.Errorf("Expected at least 1 reference, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r := addRequest.Parameter[0].References[0].Reference; l == 1 && r != userID {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, `{
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
	ok, resp, err := client.Groups.RemoveMembers(context.Background(), group, userID)
	assert.NotNil(t, resp)
	assert.Nil(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	ok, resp, err = client.Groups.RemoveMembers(context.Background(), group, "foo")
	if !assert.NotNil(t, err) {
		return
	}
	assert.NotNil(t, resp)
	assert.Nil(t, ok)
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
			body, err := io.ReadAll(r.Body)
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
			_, _ = io.WriteString(w, `{
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
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, resp)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
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
			_, _ = io.WriteString(w, `{
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
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 1, len(*roles))
	assert.Equal(t, roleID, (*roles)[0].ID)
	assert.Equal(t, "TDRALL", (*roles)[0].Name)
}

func TestAddServicesAndDevices(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	eTag := "RonSwanson"
	identityID := "f5fe538f-c3b5-4454-8774-cd3789f59b9a"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	managingOrgID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc1"
	groupName := "TestGroup"
	groupDescription := "Test Group Description"

	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.Header().Set("ETag", eTag)
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
			"name": "`+groupName+`",
			"description": "`+groupDescription+`",
			"managingOrganization": "`+managingOrgID+`",
			"id": "`+groupID+`"
			}`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$assign", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
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
			if !(addRequest.MemberType == "SERVICE" || addRequest.MemberType == "DEVICE") {
				t.Errorf("Expected SERVICE or DEVICE MemberType, got: %s", addRequest.MemberType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			basePath := "services"
			if addRequest.MemberType == "DEVICE" {
				basePath = "devices"
			}
			if l := len(addRequest.Value); l != 1 {
				t.Errorf("Expected 1 value, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Header.Get("If-Match") != eTag {
				w.WriteHeader(http.StatusBadRequest)
			}
			if n := addRequest.Value[0]; n != identityID {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"`+basePath+`/`+identityID+`"
							]
						}
					]
				}
				`)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "value": [
    "`+identityID+`"
  ],
  "type": "SERVICE"
}`)
		}
	})
	var group Group
	group.ID = groupID
	ok, resp, err := client.Groups.AddServices(context.Background(), group, identityID)
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	}
	assert.NotNil(t, ok)
	assert.Nil(t, err)
	ok, resp, err = client.Groups.AddServices(context.Background(), group, "foo")
	assert.NotNil(t, resp)
	assert.Nil(t, ok)
	assert.NotNil(t, err)
	group.ID = groupID

	ok, resp, err = client.Groups.AddDevices(context.Background(), group, identityID)
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	}
	assert.NotNil(t, ok)
	assert.Nil(t, err)
	ok, resp, err = client.Groups.AddDevices(context.Background(), group, "foo")
	assert.NotNil(t, resp)
	assert.Nil(t, ok)
	assert.NotNil(t, err)
}

func TestRemoveServicesAndDevices(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	eTag := "RonSwanson"
	identityID := "f5fe538f-c3b5-4454-8774-cd3789f59b9a"
	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	managingOrgID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc1"
	groupName := "TestGroup"
	groupDescription := "Test Group Description"
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.Header().Set("ETag", eTag)
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
			"name": "`+groupName+`",
			"description": "`+groupDescription+`",
			"managingOrganization": "`+managingOrgID+`",
			"id": "`+groupID+`"
			}`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Group/"+groupID+"/$remove", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Unexpected EOF from reading request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var removeRequest struct {
				MemberType string   `json:"memberType"`
				Value      []string `json:"value"`
			}
			err = json.Unmarshal(body, &removeRequest)
			if err != nil {
				t.Errorf("Error parsing request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !(removeRequest.MemberType == "SERVICE" || removeRequest.MemberType == "DEVICE") {
				t.Errorf("Expected SERVICE or DEVICE MemberType, got: %s", removeRequest.MemberType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			basePath := "services"
			if removeRequest.MemberType == "DEVICE" {
				basePath = "devices"
			}
			if l := len(removeRequest.Value); l != 1 {
				t.Errorf("Expected 1 value, got: %d", l)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Header.Get("If-Match") != eTag {
				w.WriteHeader(http.StatusBadRequest)
			}
			if r := removeRequest.Value[0]; r != identityID {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, `{
					"resourceType": "OperationOutcome",
					"issues": [
						{
							"diagnostic": "failed Operations",
							"code": "not-found",
							"severity": "error",
							"location": [
								"`+basePath+`/`+r+`"
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
	_, resp, err := client.Groups.RemoveServices(context.Background(), group, identityID)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	_, resp, err = client.Groups.RemoveServices(context.Background(), group, "foo")
	assert.NotNil(t, resp)
	assert.NotNil(t, err)

	_, resp, err = client.Groups.RemoveDevices(context.Background(), group, identityID)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	ok, resp, err := client.Groups.RemoveDevices(context.Background(), group, "foo")
	assert.NotNil(t, resp)
	assert.NotNil(t, err)
	assert.Nil(t, ok)
}

func TestGetSCIMGroup(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	groupID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc0"
	muxIDM.HandleFunc("/authorize/scim/v2/Groups/"+groupID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:Group",
    "urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:Group"
  ],
  "id": "72b8908f-02e5-4939-9e20-ac099ba17e5c",
  "displayName": "GroupName",
  "urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:Group": {
    "description": "Group Description",
    "organization": {
      "value": "86d2d0ac-4b12-4546-8d20-4c09a7c87d9c",
      "$ref": "https://<idm_base_path>/authorize/scim/v2/Organizations/86d2d0ac-4b12-4546-8d20-4c09a7c87d9c"
    },
    "groupMembers": {
      "schemas": [
        "urn:ietf:params:scim:api:messages:2.0:SCIMListResponse"
      ],
      "totalResults": 1,
      "startIndex": 1,
      "itemsPerPage": 1,
      "Resources": [
        {
          "schemas": [
            "urn:ietf:params:scim:schemas:core:2.0:User",
            "urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:User"
          ],
          "id": "72b8908f-02e5-4939-9e20-ac099ba17e5c",
          "userName": "wdale",
          "name": {
            "fullName": "Mr. John Jane Doe, III",
            "familyName": "Doe",
            "givenName": "John",
            "middleName": "Jane"
          },
          "displayName": "John Doe",
          "preferredLanguage": "en-US",
          "locale": "en-US",
          "active": true,
          "emails": [
            {
              "value": "john.doe@example.com",
              "primary": true
            }
          ],
          "phoneNumbers": [
            {
              "value": "555-555-4444",
              "type": "work",
              "primary": true
            }
          ],
          "urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:User": {
            "emailVerified": true,
            "phoneVerified": false,
            "organization": {
              "value": "86d2d0ac-4b12-4546-8d20-4c09a7c87d9c",
              "$ref": "https://<idm_base_path>/authorize/scim/v2/Organizations/86d2d0ac-4b12-4546-8d20-4c09a7c87d9c"
            }
          }
        }
      ]
    }
  },
  "meta": {
    "resourceType": "Group",
    "created": "2022-03-23T07:09:17.543Z",
    "lastModified": "2022-03-23T07:09:30.500Z",
    "location": "<idm_base_path>/authorize/scim/v2/Groups/72b8908f-02e5-4939-9e20-ac099ba17e5c",
    "version": "W/\"f250dd84f0671c3\""
  }
}`)
		}
	})
	var group Group
	group.ID = groupID
	members := "USER"
	scimGroup, resp, err := client.Groups.SCIMGetGroupByID(groupID, &SCIMGetGroupOptions{
		IncludeGroupMembersType: &members,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, scimGroup) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
