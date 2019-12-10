package credentials

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// PolicyService provides operations on S3 Credentials Policies
type PolicyService struct {
	client *Client

	validate *validator.Validate
}

// Constants
const (
	CredentialsAPIVersion = "2"
)

// GetPolicyOptions describes the fields on which you can search for policies
type GetPolicyOptions struct {
	ManagingOrg *string `url:"managing-org,omitempty"`
	GroupName   *string `url:"group-name,omitempty"`
	ID          *int    `url:"id,omitempty"`
	ProductKey  *string `url:"-"`
}

// GetPolicy searches for polices
func (c *PolicyService) GetPolicy(opt *GetPolicyOptions, options ...OptionFunc) ([]*Policy, *Response, error) {

	req, err := c.client.NewRequest("GET", "core/credentials/Policy", opt, options)
	if err != nil {
		return nil, nil, err
	}
	if opt.ProductKey == nil {
		return nil, nil, ErrMissingProductKey
	}

	req.Header.Set("Api-Version", CredentialsAPIVersion)
	req.Header.Set("X-Product-Key", *opt.ProductKey)

	var policyGetResponse []*Policy

	resp, err := c.client.Do(req, &policyGetResponse)
	if err != nil {
		return nil, resp, err
	}
	// Set ProductKey
	for _, policy := range policyGetResponse {
		policy.ProductKey = *opt.ProductKey
	}
	return policyGetResponse, resp, err
}

// CreatePolicy creates a new policy for S3 Credentials
func (c *PolicyService) CreatePolicy(policy Policy) (*Policy, *Response, error) {
	if err := c.validate.Struct(policy); err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest("POST", "core/credentials/Policy", &policy, nil)
	if err != nil {
		return nil, nil, err
	}
	if policy.ProductKey == "" {
		return nil, nil, ErrMissingProductKey
	}
	req.Header.Set("Api-Version", CredentialsAPIVersion)
	req.Header.Set("X-Product-Key", policy.ProductKey)

	var createdPolicy Policy
	resp, err := c.client.Do(req, &createdPolicy)
	if err != nil {
		return nil, resp, err
	}
	if resp == nil {
		return nil, resp, ErrOperationFailed
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, resp, err
	}
	createdPolicy.ProductKey = policy.ProductKey

	return &createdPolicy, resp, nil
}

// DeleteGroup deletes the given Group
func (c *PolicyService) DeletePolicy(policy Policy) (bool, *Response, error) {
	req, err := c.client.NewRequest("DELETE", "core/credentials/Policy/"+policy.StringID(), nil, nil)
	if err != nil {
		return false, nil, err
	}
	if policy.ProductKey == "" {
		return false, nil, ErrMissingProductKey
	}
	req.Header.Set("api-version", CredentialsAPIVersion)
	req.Header.Set("X-Product-Key", policy.ProductKey)

	var deleteResponse interface{}

	resp, err := c.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}

	return true, resp, nil
}
