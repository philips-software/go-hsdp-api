package iam

import "time"

// Group represents an IAM group resource
type Group struct {
	ID                   string `json:"id,omitempty" validate:""`
	Name                 string `json:"name,omitempty" validate:"required"`
	Description          string `json:"description,omitempty" validate:""`
	ManagingOrganization string `json:"managingOrganization,omitempty" validate:"required"`
}

// GroupResource is the resource response of a Group search operation
type GroupResource struct {
	ID               string `json:"_id"`
	ResourceType     string `json:"resourceType"`
	GroupName        string `json:"groupName"`
	OrgID            string `json:"orgId"`
	GroupDescription string `json:"groupDescription"`
}

// SCIMGroup is the resource returned when getting group resources
type SCIMGroup struct {
	Schemas        []string       `json:"schemas"`
	ID             string         `json:"id"`
	DisplayName    string         `json:"displayName"`
	ExtensionGroup ExtensionGroup `json:"urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:Group"`
	Meta           *Meta          `json:"meta,omitempty"`
}

type ExtensionGroup struct {
	Description  string           `json:"description"`
	Organization Attribute        `json:"organization"`
	GroupMembers SCIMListResponse `json:"groupMembers"`
}

type ExtensionUser struct {
	EmailVerified bool      `json:"emailVerified"`
	PhoneVerified bool      `json:"phoneVerified"`
	Organization  Attribute `json:"organization"`
}

type SCIMName struct {
	FullName   string `json:"fullName,omitempty"`
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
}

type SCIMListResponse struct {
	Schemas      []string           `json:"schemas"`
	TotalResults int                `json:"totalResults"`
	StartIndex   int                `json:"startIndex"`
	ItemsPerPage int                `json:"itemsPerPage"`
	Resources    []SCIMListResource `json:"Resources"`
}

type SCIMListResource struct {
	Schemas      []string  `json:"schemas"`
	ID           string    `json:"id"`
	Organization Attribute `json:"organization,omitempty"`
	Application  Attribute `json:"application,omitempty"`
	Active       bool      `json:"active,omitempty"`
	SCIMCoreUser
	SCIMService
	SCIMDevice
	ExtensionUser ExtensionUser `json:"urn:ietf:params:scim:schemas:extension:philips:hsdp:2.0:User,omitempty"`
}

type SCIMService struct {
	ServiceId string     `json:"serviceId,omitempty"`
	ExpiresOn *time.Time `json:"expiresOn,omitempty"`
	// Organization
	// Application
}

type SCIMDevice struct {
	LoginID string `json:"loginId,omitempty"`
	// Organization
	// Application
	// Active
}

type SCIMCoreUser struct {
	UserName          string      `json:"userName,omitempty"`
	DisplayName       string      `json:"displayName,omitempty"`
	Name              SCIMName    `json:"name,omitempty"`
	PreferredLanguage string      `json:"preferredLanguage,omitempty"`
	Locale            string      `json:"locale,omitempty"`
	Emails            []Attribute `json:"emails,omitempty"`
	PhoneNumbers      []Attribute `json:"phoneNumbers,omitempty"`
}

const (
	GroupMemberTypeUser    = "USER"
	GroupMemberTypeDevice  = "DEVICE"
	GroupMemberTypeService = "SERVICE"
)
