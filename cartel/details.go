package cartel

import (
	"encoding/json"
)

type LdapGroups []string

func (lg *LdapGroups) UnmarshalJSON(b []byte) error {
	var a []string
	if b[0] == '[' { // String array
		err := json.Unmarshal(b, &a)
		if err != nil {
			return err
		}
		*lg = a
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*lg = []string{s}
	return nil
}

type InstanceDetails struct {
	BlockDevices   []string          `json:"block_devices,omitempty"`
	InstanceID     string            `json:"instance_id"`
	InstanceType   string            `json:"instance_type,omitempty"`
	LaunchTime     string            `json:"launch_time,omitempty"`
	LdapGroups     LdapGroups        `json:"ldap_groups,omitempty"`
	PrivateAddress string            `json:"private_address,omitempty"`
	Protection     bool              `json:"protection,omitempty"`
	PublicAddress  string            `json:"public_address,omitempty"`
	Role           string            `json:"role"`
	SecurityGroups []string          `json:"security_groups,omitempty"`
	State          string            `json:"state,omitempty"`
	Subnet         string            `json:"subnet,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
	Vpc            string            `json:"vpc,omitempty"`
	Zone           string            `json:"zone,omitempty"`
	Owner          string            `json:"owner,omitempty"`
	NameTag        string            `json:"name_tag,omitempty"`
}

type DetailsResponse map[string]InstanceDetails

func (c *Client) GetDetailsMulti(tags ...string) (*DetailsResponse, *Response, error) {
	var body RequestBody
	body.NameTag = tags

	req, err := c.newRequest("POST", "v3/api/instance_details", &body, nil)
	if err != nil {
		return nil, nil, err
	}

	var detailResponse []map[string]InstanceDetails

	resp, err := c.do(req, &detailResponse)
	response := make(DetailsResponse, len(detailResponse))
	for _, r := range detailResponse {
		for k, v := range r {
			response[k] = v
		}
	}
	return &response, resp, err
}

func (c *Client) GetDetails(tag string) (*InstanceDetails, *Response, error) {
	details, resp, err := c.GetDetailsMulti(tag)
	if err != nil {
		return nil, resp, err
	}
	if len(*details) == 0 {
		return nil, resp, ErrNotFound
	}
	id := (*details)[tag]
	return &id, resp, nil
}
