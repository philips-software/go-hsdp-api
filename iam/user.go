package iam

import (
	"time"
)

// User represents a user profile in IAM
type User struct {
	PreferredLanguage             string             `json:"preferredLanguage"`
	PreferredCommunicationChannel string             `json:"preferredCommunicationChannel,omitempty"`
	EmailAddress                  string             `json:"emailAddress"`
	PhoneNumber                   string             `json:"phoneNumber,omitempty"`
	ID                            string             `json:"id"`
	LoginID                       string             `json:"loginId"`
	Name                          Name               `json:"name"`
	ManagingOrganization          string             `json:"managingOrganization"`
	PasswordStatus                UserPasswordStatus `json:"passwordStatus"`
	Memberships                   []UserMembership   `json:"memberships,omitempty"`
	AccountStatus                 UserAccountStatus  `json:"accountStatus"`
	ConsentedApps                 []string           `json:"consentedApps,omitempty"`
	Delegations                   UserDelegation     `json:"delegations,omitempty"`
}

type UserDelegation struct {
	Granted  []UserDelegator `json:"granted"`
	Received []UserDelegator `json:"received"`
}

type UserDelegator struct {
	DelegateeID string `json:"delegateeId"`
	ValidFrom   string `json:"validFrom"`
	ValidUntil  string `json:"validUntil"`
}

type UserMembership struct {
	OrganizationID   string   `json:"organizationId"`
	OrganizationName string   `json:"organizationName"`
	Roles            []string `json:"roles"`
	Groups           []string `json:"groups"`
}

type UserAccountStatus struct {
	LastLoginTime          time.Time `json:"lastLoginTime"`
	MfaStatus              string    `json:"mfaStatus"`
	EmailVerified          bool      `json:"emailVerified"`
	Disabled               bool      `json:"disabled"`
	AccountLockedOn        time.Time `json:"accountLockedOn"`
	AccountLockedUntil     time.Time `json:"accountLockedUntil"`
	NumberOfInvalidAttempt int       `json:"numberOfInvalidAttempt"`
	LastInvalidAttemptedOn time.Time `json:"lastInvalidAttemptedOn"`
}

type UserPasswordStatus struct {
	PasswordExpiresOn time.Time `json:"passwordExpiresOn"`
	PasswordChangedOn time.Time `json:"passwordChangedOn"`
}

// Person represents an IAM user resource
type Person struct {
	ID string `json:"id,omitempty" validate:"omitempty"`
	// Pattern: ^((?![~`!#%^&*()+={}[\\]|/\\\\<>,;:\"'?])[\\S])*$
	LoginID                       string         `json:"loginId" validate:"required"`
	ResourceType                  string         `json:"resourceType,omitempty" validate:"required" enum:"Person"`
	Name                          Name           `json:"name" validate:"required"`
	Telecom                       []TelecomEntry `json:"telecom,omitempty" validate:"min=1"`
	Address                       []AddressEntry `json:"address,omitempty"`
	Description                   string         `json:"description,omitempty"`
	ManagingOrganization          string         `json:"managingOrganization,omitempty"`
	PreferredLanguage             string         `json:"preferredLanguage,omitempty"`
	PreferredCommunicationChannel string         `json:"preferredCommunicationChannel,omitempty"`
	IsAgeValidated                string         `json:"isAgeValidated,omitempty"`
	Password                      string         `json:"password,omitempty"`
	Disabled                      bool           `json:"disabled"`
	Loaded                        bool           `json:"-"`
}

// Contact describes contact details of a Profile
type Contact struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	MobilePhone  string `json:"mobilePhone,omitempty"`
	WorkPhone    string `json:"workPhone,omitempty"`
	HomePhone    string `json:"homePhone,omitempty"`
}

// Address describes an address of a Profile
type Address struct {
	Use        string   `json:"use,omitempty" enum:"home|work|temp|old"`
	Text       string   `json:"text,omitempty"`
	City       string   `json:"city,omitempty"`
	State      string   `json:"state,omitempty"`
	Line       []string `json:"line,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
	Building   string   `json:"building,omitempty"`
	Street     string   `json:"street,omitempty"`
	IsPrimary  string   `json:"isPrimary,omitempty" enum:"yes|no"`
}

func (a *Address) IsBlank() bool {
	return len(a.Use) == 0 && len(a.Text) == 0 && len(a.City) == 0 && len(a.State) == 0 &&
		len(a.Line) == 0 && len(a.PostalCode) == 0 && len(a.Country) == 0 && len(a.Building) == 0 &&
		len(a.Street) == 0 && len(a.IsPrimary) == 0
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
	ID                string     `json:"-"`
	GivenName         string     `json:"givenName"`
	MiddleName        string     `json:"middleName"`
	FamilyName        string     `json:"familyName"`
	BirthDay          *time.Time `json:"birthDay,omitempty"`
	DisplayName       string     `json:"displayName,omitempty"`
	Gender            string     `json:"gender,omitempty" enum:"Male|Female"`
	Country           string     `json:"country,omitempty"`
	Contact           Contact    `json:"contact,omitempty"`
	Addresses         []Address  `json:"addresses,omitempty"`
	PreferredLanguage string     `json:"preferredLanguage,omitempty"`
}

// PruneBlankAddresses removes addresses which are blank
func (p *Profile) PruneBlankAddresses() {
	if len(p.Addresses) == 0 {
		return
	}
	pruned := make([]Address, 0)
	for _, a := range p.Addresses {
		if !a.IsBlank() {
			pruned = append(pruned, a)
		}
	}
	if len(pruned) == len(p.Addresses) {
		return
	}
	p.Addresses = pruned
}

// MergeUser merges User into legacy Profile
func (p *Profile) MergeUser(user *User) {
	p.GivenName = user.Name.Given
	p.FamilyName = user.Name.Family
	// See INC0058741 for backround for this workaround
	if p.MiddleName == "" {
		p.MiddleName = " "
	}
	p.ID = user.ID
	p.Contact.EmailAddress = user.EmailAddress
	p.PreferredLanguage = user.PreferredLanguage
}

// Name entity
type Name struct {
	Text   string `json:"text,omitempty"`
	Family string `json:"family" validate:"required"`
	Given  string `json:"given" validate:"required"`
	Prefix string `json:"prefix,omitempty"`
}

// TelecomEntry entity
type TelecomEntry struct {
	System string `json:"system" enum:"mobile|fax|email|url"`
	Value  string `json:"value"`
}

// AddressEntry entity
type AddressEntry struct {
	Use        string   `json:"use,omitempty"`
	Text       string   `json:"text,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       string   `json:"city,omitempty"`
	State      string   `json:"state,omitempty"`
	Country    string   `json:"country,omitempty"`
	Postalcode string   `json:"postalcode,omitempty"`
}
