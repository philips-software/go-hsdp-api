package iam

import (
	"bytes"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	PasswordPolicyAPIVersion = "1"
)

type PasswordPoliciesService struct {
	client   *Client
	validate *validator.Validate
}

type ChallengePolicy struct {
	DefaultQuestions     []string `json:"defaultQuestions"`
	MinQuestionCount     int      `json:"minQuestionCount"`
	MinAnswerCount       int      `json:"minAnswerCount"`
	MaxIncorrectAttempts int      `json:"maxIncorrectAttempts"`
}

type PasswordPolicy struct {
	ID                   string `json:"id,omitempty"`
	ManagingOrganization string `json:"managingOrganization"`
	ExpiryPeriodInDays   int    `json:"expiryPeriodInDays"`
	HistoryCount         int    `json:"historyCount"`
	Complexity           struct {
		MinLength       int `json:"minLength"`
		MaxLength       int `json:"maxLength"`
		MinNumerics     int `json:"minNumerics"`
		MinUpperCase    int `json:"minUpperCase"`
		MinLowerCase    int `json:"minLowerCase"`
		MinSpecialChars int `json:"minSpecialChars"`
	} `json:"complexity"`
	ChallengesEnabled bool             `json:"challengesEnabled"`
	ChallengePolicy   *ChallengePolicy `json:"challengePolicy,omitempty"`
	Meta              *Meta            `json:"meta,omitempty"`
}

// GetPasswordPolicyByID retrieves a Password policy by ID
func (p *PasswordPoliciesService) GetPasswordPolicyByID(id string) (*PasswordPolicy, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/PasswordPolicy/"+id, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", PasswordPolicyAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var policy PasswordPolicy

	resp, err := p.client.Do(req, &policy)
	if err != nil {
		return nil, resp, err
	}
	if policy.ID != id {
		return nil, resp, ErrNotFound
	}
	return &policy, resp, err
}

// UpdatePasswordPolicy updates a password policy
func (p *PasswordPoliciesService) UpdatePasswordPolicy(policy *PasswordPolicy) (*PasswordPolicy, *Response, error) {

	req, _ := p.client.NewRequest(IDM, "PUT", "authorize/identity/PasswordPolicy/"+policy.ID, policy, nil)
	req.Header.Set("api-version", PasswordPolicyAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	if policy.Meta == nil {
		return nil, nil, ErrMissingEtagInformation
	}
	req.Header.Set("If-Match", policy.Meta.Version)

	var updatedPolicy PasswordPolicy
	resp, err := p.client.Do(req, &updatedPolicy)

	if err != nil {
		return nil, resp, err
	}
	return &updatedPolicy, resp, nil

}

// CreatePasswordPolicy creates a password policy
func (p *PasswordPoliciesService) CreatePasswordPolicy(policy PasswordPolicy) (*PasswordPolicy, *Response, error) {
	if err := p.validate.Struct(policy); err != nil {
		return nil, nil, err
	}
	req, _ := p.client.NewRequest(IDM, "POST", "authorize/identity/PasswordPolicy", &policy, nil)
	req.Header.Set("api-version", PasswordPolicyAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var createdPolicy PasswordPolicy

	resp, err := p.client.Do(req, &createdPolicy)
	if err != nil {
		return nil, resp, err
	}
	return &createdPolicy, resp, err
}

// DeletePasswordPolicy deletes the given password policy
func (p *PasswordPoliciesService) DeletePasswordPolicy(policy PasswordPolicy) (bool, *Response, error) {
	req, err := p.client.NewRequest(IDM, "DELETE", "authorize/identity/PasswordPolicy/"+policy.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", PasswordPolicyAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var deleteResponse bytes.Buffer

	resp, err := p.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

// GetPasswordPolicyOptions describes the criteria for looking up password polices
type GetPasswordPolicyOptions struct {
	OrganizationID *string `url:"organizationId,omitempty"`
}

// GetPasswordPolicies looks up clients based on GetClientsOptions
func (p *PasswordPoliciesService) GetPasswordPolicies(opt *GetPasswordPolicyOptions, options ...OptionFunc) (*[]PasswordPolicy, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/PasswordPolicy", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", PasswordPolicyAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	var bundleResponse struct {
		Total int              `json:"total"`
		Entry []PasswordPolicy `json:"entry"`
	}

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	return &bundleResponse.Entry, resp, err
}
