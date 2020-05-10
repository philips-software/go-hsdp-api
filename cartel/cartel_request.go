package cartel

import "encoding/json"

type CartelRequestBody struct {
	Token         string            `json:"token,omitempty"`
	NameTag       []string          `json:"name-tag,omitempty"`
	Role          string            `json:"role,omitempty"`
	SecurityGroup []string          `json:"security_group,omitempty"`
	Image         string            `json:"image,omitempty"`
	LDAPGroups    []string          `json:"ldap_groups,omitempty"`
	ExtraScripts  []string          `json:"extra_scripts,omitempty"`
	InstanceType  string            `json:"instance_type,omitempty"`
	NumVolumes    int               `json:"num_vols,omitempty"`
	VolSize       int               `json:"vol_size,omitempty"`
	VolumeType    string            `json:"vol_type,omitempty"`
	IOPs          int               `json:"iops,omitempty"`
	EncryptVols   bool              `json:"encrypt_vols,omitempty"`
	SubnetType    string            `json:"subnet_type,omitempty"`
	Subnet        string            `json:"subnet,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Protect       bool              `json:"protect"`
	VpcId         string            `json:"vpc_id,omitempty"`
}

func (crb *CartelRequestBody) ToJson() []byte {
	req, _ := json.Marshal(crb)
	return req
}
