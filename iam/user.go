package iam

import "time"

// User represents a user profile in IAM
type User struct {
	PreferredLanguage string `json:"preferredLanguage"`
	EmailAddress      string `json:"emailAddress"`
	ID                string `json:"id"`
	LoginID           string `json:"loginId"`
	Name              struct {
		Given  string `json:"given"`
		Family string `json:"family"`
	} `json:"name"`
	ManagingOrganization string `json:"managingOrganization"`
	PasswordStatus       struct {
		PasswordExpiresOn time.Time `json:"passwordExpiresOn"`
		PasswordChangedOn time.Time `json:"passwordChangedOn"`
	} `json:"passwordStatus"`
	Memberships []struct {
		OrganizationID   string   `json:"organizationId"`
		OrganizationName string   `json:"organizationName"`
		Roles            []string `json:"roles"`
		Groups           []string `json:"groups"`
	} `json:"memberships"`
	AccountStatus struct {
		LastLoginTime          time.Time `json:"lastLoginTime"`
		MfaStatus              string    `json:"mfaStatus"`
		EmailVerified          bool      `json:"emailVerified"`
		Disabled               bool      `json:"disabled"`
		AccountLockedOn        time.Time `json:"accountLockedOn"`
		AccountLockedUntil     time.Time `json:"accountLockedUntil"`
		NumberOfInvalidAttempt int       `json:"numberOfInvalidAttempt"`
		LastInvalidAttemptedOn time.Time `json:"lastInvalidAttemptedOn"`
	} `json:"accountStatus"`
	ConsentedApps []string `json:"consentedApps"`
}

// Person represents an IAM user resource
type Person struct {
	ID string `json:"id,omitempty" validate:"omitempty"`
	// Pattern: ^((?![~`!#%^&*()+={}[\\]|/\\\\<>,;:\"'?])[\\S])*$
	LoginID              string         `json:"loginId" validate:"required"`
	ResourceType         string         `json:"resourceType,omitempty" validate:"required" enum:"Person"`
	Name                 Name           `json:"name" validate:"required"`
	Telecom              []TelecomEntry `json:"telecom,omitempty" validate:"min=1"`
	Address              []AddressEntry `json:"address,omitempty"`
	Description          string         `json:"description,omitempty"`
	ManagingOrganization string         `json:"managingOrganization,omitempty"`
	PreferredLanguage    string         `json:"preferredLanguage,omitempty"`
	IsAgeValidated       string         `json:"isAgeValidated,omitempty"`
	Disabled             bool           `json:"disabled"`
	Loaded               bool           `json:"-"`
}

// Contact describes contact details of a Profile
type Contact struct {
	EmailAddress string
	MobilePhone  string
	WorkPhone    string
	HomePhone    string
}

// Address describes an address of a Profile
type Address struct {
	Use        string `json:"use" enum:"home|work|temp|old"`
	Text       string `json:"text"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Building   string `json:"building"`
	Street     string `json:"street"`
}

// Period defines a given time period for use in Profile context
type Period struct {
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
}

// Profile describes the response from legacy User APIs
// The response does not correspond to the object that is used to create a user
type Profile struct {
	GivenName         string    `json:"givenName"`
	MiddleName        string    `json:"middleName"`
	FamilyName        string    `json:"familyName"`
	BirthDay          string    `json:"birthDay"`
	DisplayName       string    `json:"displayName"`
	Gender            string    `json:"gender" enum:"Male|Female"`
	Country           string    `json:"country"`
	Addresses         []Address `json:"addresses"`
	PreferredLanguage string    `json:"preferredLanguage"`
}

// Name entity
type Name struct {
	Text   string `json:"text"`
	Family string `json:"family" validate:"required"`
	Given  string `json:"given" validate:"required"`
	Prefix string `json:"prefix"`
}

// TelecomEntry entity
type TelecomEntry struct {
	System string `json:"system" enum:"mobile|fax|email|url"`
	Value  string `json:"value"`
}

// AddressEntry entity
type AddressEntry struct {
	Use        string `json:"use"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Postalcode string `json:"postalcode"`
}
