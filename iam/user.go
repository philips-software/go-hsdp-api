package iam

type User struct {
	ID                   string         `json:"id,omitempty"`
	ResourceType         string         `json:"resourceType,omityempty"`
	Name                 Name           `json:"name"`
	Telecom              []TelecomEntry `json:"telecom,omitempty"`
	Address              []AddressEntry `json:"address,omitempty"`
	Description          string         `json:"description,omitempty"`
	ManagingOrganization string         `json:"managingOrganization,omitempty"`
	PreferredLanguage    string         `json:"preferredLanguage,omitempty"`
	IsAgeValidated       string         `json:"isAgeValidated,omitempty"`
}

type Name struct {
	Text   string `json:"text"`
	Family string `json:"family"`
	Given  string `json:"given"`
	Prefix string `json:"prefix"`
}

type TelecomEntry struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type AddressEntry struct {
	Use        string `json:"use"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Postalcode string `json:"postalcode"`
}

func (u *User) Validate() error {
	return nil
}
