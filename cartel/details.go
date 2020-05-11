package cartel

import "time"

type InstanceDetails struct {
	BlockDevices   []string    `json:"block_devices"`
	InstanceID     string      `json:"instance_id"`
	InstanceType   string      `json:"instance_type"`
	LaunchTime     time.Time   `json:"launch_time"`
	LdapGroups     []string    `json:"ldap_groups"`
	PrivateAddress string      `json:"private_address"`
	Protection     bool        `json:"protection"`
	PublicAddress  interface{} `json:"public_address"`
	Role           string      `json:"role"`
	SecurityGroups []string    `json:"security_groups"`
	State          string      `json:"state"`
	Subnet         string      `json:"subnet"`
	Tags           struct {
		Billing string `json:"billing"`
	} `json:"tags"`
	Vpc  string `json:"vpc"`
	Zone string `json:"zone"`
}

type DetailsResponse map[string]InstanceDetails

func (c *Client) GetDetails(tags ...string) (*DetailsResponse, *Response, error) {
	var body CartelRequestBody
	body.NameTag = tags

	req, err := c.NewRequest("POST", "v3/api/instance_details", body, nil)
	if err != nil {
		return nil, nil, err
	}

	var detailResponse []map[string]InstanceDetails

	resp, err := c.Do(req, &detailResponse)
	var response DetailsResponse
	for _, r := range detailResponse {
		for k, v := range r {
			response[k] = v
		}
	}
	return &response, resp, err
}

func (c *Client) GetDetail(tag string) (*InstanceDetails, *Response, error) {
	details, resp, err := c.GetDetails(tag)
	if err != nil {
		return nil, resp, err
	}
	if len(*details) == 0 {
		return nil, resp, ErrNotFound
	}
	id := (*details)[tag]
	return &id, resp, nil
}
