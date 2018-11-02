package iam

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeffail/gabs"
)

const (
	userAPIVersion = "1"
)

// GetUserOptions describes search criteria for looking up users
type GetUserOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"Id,omitempty"`
	Name           *string `url:"name,omitempty"`
	LoginID        *string `url:"loginId,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	PageSize       *string `url:"pageSize,omitempty"`
	PageNumber     *string `url:"pageNumber,omitempty"`
}

// UsersService provides operations on IAM User resources
type UsersService struct {
	client *Client
}

// Parameters holds parameters
type Parameters struct {
	ResourceType string  `json:"resourceType"`
	Parameter    []Param `json:"parameter"`
}

// Param describes a resource
type Param struct {
	Name     string   `json:"name"`
	Resource Resource `json:"resource"`
}

// Resource holds a resource
type Resource struct {
	LoginID          string `json:"loginId,omitempty"`
	ConfirmationCode string `json:"confirmationCode,omitempty"`
	OldPassword      string `json:"oldPassword,omitempty"`
	NewPassword      string `json:"newPassword,omitempty"`
	Context          string `json:"context,omitempty"`
}

// UserList holds a paginated lists of users
type UserList struct {
	Users       []Person
	PageNumber  int
	PageSize    int
	HasNextPage bool
}

// CreateUser creates a new IAM user.
func (u *UsersService) CreateUser(firstName, lastName, emailID, phoneNumber, organizationID string) (bool, *Response, error) {
	person := &Person{
		ResourceType: "Person",
		Name: Name{
			Family: lastName,
			Given:  firstName,
		},
		Telecom: []TelecomEntry{
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
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User", person, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}
	var doFunc func(*http.Request, interface{}) (*Response, error)

	if organizationID == "" { // Self registration
		doFunc = u.client.DoSigned
	} else { // Admin registration
		doFunc = u.client.Do
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	return ok, resp, err
}

// RecoverPassword triggers the recovery flow for the given user
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
	return u.userAction(body, "$recover-password")
}

// ResendActivation resends an activation email to the given user
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
	return u.userAction(body, "$resend-activation")
}

func (u *UsersService) userAction(body *Parameters, action string) (bool, *Response, error) {
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+action, body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted)
	return ok, resp, err
}

// SetPassword sets the password of a user given a correct confirmation code
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
	return u.userAction(body, "$set-password")
}

// ChangePassword changes the password. The current pasword must be provided as well.
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
	return u.userAction(body, "$change-password")
}

// GetUsers looks up users by search criteria specified in GetUserOptions
func (u *UsersService) GetUsers(opts *GetUserOptions, options ...OptionFunc) (*UserList, *Response, error) {
	req, err := u.client.NewRequest(IDM, "GET", "security/users", opts, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)
	var bundleResponse struct {
		Exchange struct {
			Users []struct {
				UserUUID string `json:"userUUID"`
			}
			NextPageExists bool `json:"nextPageExists"`
		}
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var list UserList

	list.HasNextPage = bundleResponse.Exchange.NextPageExists
	list.Users = make([]Person, len(bundleResponse.Exchange.Users))
	for i, u := range bundleResponse.Exchange.Users {
		list.Users[i] = Person{ID: u.UserUUID}
	}

	return &list, resp, err
}

// GetUserByID looks up a user by UUID
func (u *UsersService) GetUserByID(uuid string) (*Person, *Response, error) {
	req, _ := u.client.NewRequest(IDM, "GET", "security/users/"+uuid, nil, nil)
	var user interface{}

	resp, err := u.client.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	m, err := json.Marshal(user)
	if err != nil {
		return nil, resp, fmt.Errorf("error parsing json response")
	}
	jsonParsed, err := gabs.ParseJSON(m)
	if statusCode, ok := jsonParsed.Path("responseCode").Data().(string); !ok || statusCode != "200" {
		return nil, resp, fmt.Errorf("responseCode: %s", statusCode)
	}
	email, ok := jsonParsed.Path("exchange.loginId").Data().(string)
	if !ok {
		return nil, resp, fmt.Errorf("Invalid response")
	}
	r := jsonParsed.Path("exchange.profile")
	first := r.Path("givenName").Data().(string)
	last := r.Path("familyName").Data().(string)
	disabled := r.Path("disabled").Data().(bool)
	// TODO use Profile here
	var foundUser Person
	foundUser.Name.Family = last
	foundUser.Name.Given = first
	foundUser.Disabled = disabled
	foundUser.Telecom = append(foundUser.Telecom, TelecomEntry{
		System: "email",
		Value:  email,
	})
	return &foundUser, resp, nil
}

// GetUserIDByLoginID looks up the UUID of a user by LoginID (email address)
func (u *UsersService) GetUserIDByLoginID(loginID string) (string, *Response, error) {
	req, err := u.client.NewRequest(IDM, "GET", "security/users?loginId="+loginID, nil, nil)
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

// SetMFA activate Multi-Factor-Authentication for the given UUID. See also SetMFAByLoginID.
func (u *UsersService) SetMFA(userID string, activate bool) (bool, *Response, error) {
	activateString := "true"
	if !activate {
		activateString = "false"
	}
	body := &struct {
		Activate string `json:"activate"`
	}{activateString}
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+userID+"/$mfa", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.Do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode == http.StatusAccepted)
	return ok, resp, nil
}

// SetMFAByLoginID enabled Multi-Factor-Authentication for the given user. Only OrgAdmins can do this.
func (u *UsersService) SetMFAByLoginID(loginID string, activate bool) (bool, *Response, error) {
	userUUID, _, err := u.GetUserIDByLoginID(loginID)
	if err != nil {
		return false, nil, err
	}
	return u.SetMFA(userUUID, activate)
}
