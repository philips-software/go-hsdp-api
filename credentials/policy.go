package credentials

import "strconv"

type Policy struct {
	ID           int    `json:"id,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Conditions   struct {
		ManagingOrganizations []string `json:"managingOrganizations"`
		Groups                []string `json:"groups"`
	} `json:"conditions"`
	Allowed struct {
		Resources []string `json:"resources"`
		Actions   []string `json:"actions"`
	} `json:"allowed"`

	ProductKey string `json:"-"`
}

func (p Policy) StringID() string {
	return strconv.Itoa(p.ID)
}
