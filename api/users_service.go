package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hsdp/go-hsdp-iam/iam"
	"github.com/jeffail/gabs"
)

const (
	UserAPIVersion = "1"
)

type GetUserOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"Id,omitempty"`
	Name           *string `url:"name,omitempty"`
}

type UsersService struct {
	client *Client
}

type Parameters struct {
	ResourceType string  `json:"resourceType"`
	Parameter    []Param `json:"parameter"`
}

type Param struct {
	Name     string   `json:"name"`
	Resource Resource `json:"resource"`
}

type Resource struct {
	LoginID          string `json:"loginId,omitempty"`
	ConfirmationCode string `json:"confirmationCode,omitempty"`
	OldPassword      string `json:"oldPassword,omitempty"`
	NewPassword      string `json:"newPassword,omitempty"`
	Context          string `json:"context,omitempty"`
}

func (u *UsersService) CreateUser(firstName, lastName, emailID, phoneNumber, organizationID string) (bool, *Response, error) {
	person := &iam.User{
		ResourceType: "Person",
		Name: iam.Name{
			Family: lastName,
			Given:  firstName,
		},
		Telecom: []iam.TelecomEntry{
			{
				System: "mobile",
				Value:  phoneNumber,
			},
			{
				System: "email",
				Value:  emailID,
			},
		},
		ManagingOrganization: organizationID,
		IsAgeValidated:       "true",
	}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User", person, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	return ok, resp, err
}

func (u *UsersService) RecoverPassword(loginID string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "recoverPassword",
				Resource: Resource{
					LoginID: loginID,
				},
			},
		},
	}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User/$recover-password", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)

	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && resp.StatusCode == http.StatusOK
	return ok, resp, err
}

func (u *UsersService) ResendActivation(loginID string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "resendOTP",
				Resource: Resource{
					LoginID: loginID,
				},
			},
		},
	}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User/$resend-activation", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted)
	return ok, resp, err
}

func (u *UsersService) SetPassword(loginID, confirmationCode, newPassword, context string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "setPassword",
				Resource: Resource{
					LoginID:          loginID,
					ConfirmationCode: confirmationCode,
					NewPassword:      newPassword,
					Context:          context,
				},
			},
		},
	}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User/$set-password", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && resp.StatusCode == http.StatusOK
	return ok, resp, err
}

func (u *UsersService) ChangePassword(loginID, oldPassword, newPassword string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "changePassword",
				Resource: Resource{
					LoginID:     loginID,
					OldPassword: oldPassword,
					NewPassword: newPassword,
				},
			},
		},
	}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User/$change-password", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && resp.StatusCode == http.StatusOK
	return ok, resp, err
}

func (u *UsersService) GetUserIDByLoginID(loginID string) (string, *Response, error) {
	req, err := u.client.NewIDMRequest("GET", "security/users?loginId="+loginID, nil, nil)
	var d interface{}

	resp, err := u.client.Do(req, &d)
	if err != nil {
		return "", resp, err
	}
	m, err := json.Marshal(d)
	if err != nil {
		return "", resp, fmt.Errorf("error parsing json response")
	}
	jsonParsed, err := gabs.ParseJSON(m)
	if statusCode, ok := jsonParsed.Path("responseCode").Data().(string); !ok || statusCode != "200" {
		return "", resp, fmt.Errorf("responseCode: %s", statusCode)
	}

	r := jsonParsed.Path("exchange.users").Index(0)
	userUUID, ok := r.Path("userUUID").Data().(string)
	if !ok {
		return "", resp, fmt.Errorf("lookup failed")
	}
	return userUUID, resp, nil

}

func (u *UsersService) SetMFA(userID string, activate bool) (bool, *Response, error) {
	activateString := "true"
	if !activate {
		activateString = "false"
	}
	body := &struct {
		Activate string `json:"activate"`
	}{activateString}
	req, err := u.client.NewIDMRequest("POST", "authorize/identity/User/"+userID+"/$mfa", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", UserAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.Do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode == http.StatusAccepted)
	return ok, resp, nil
}

func (u *UsersService) SetMFAByLoginID(loginID string, activate bool) (bool, *Response, error) {
	userUUID, _, err := u.GetUserIDByLoginID(loginID)
	if err != nil {
		return false, nil, err
	}
	return u.SetMFA(userUUID, activate)
}
