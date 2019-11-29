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

func (p *Policy) StringID() string {
	return strconv.Itoa(p.ID)
}

// Equals determines of other Policy is equavalent
func (p *Policy) Equals(other *Policy) bool {
	if p.ID != other.ID {
		return false
	}
	if p.ResourceType != other.ResourceType {
		return false
	}
	if len(difference(p.Conditions.ManagingOrganizations, other.Conditions.ManagingOrganizations)) > 0 {
		return false
	}
	if len(difference(other.Conditions.ManagingOrganizations, p.Conditions.ManagingOrganizations)) > 0 {
		return false
	}
	if len(difference(p.Conditions.Groups, other.Conditions.Groups)) > 0 {
		return false
	}
	if len(difference(other.Conditions.Groups, p.Conditions.Groups)) > 0 {
		return false
	}
	if len(difference(p.Allowed.Resources, other.Allowed.Resources)) > 0 {
		return false
	}
	if len(difference(other.Allowed.Resources, p.Allowed.Resources)) > 0 {
		return false
	}
	if len(difference(p.Allowed.Actions, other.Allowed.Actions)) > 0 {
		return false
	}
	if len(difference(other.Allowed.Actions, p.Allowed.Actions)) > 0 {
		return false
	}
	if p.ProductKey != other.ProductKey {
		return false
	}
	return true
}
