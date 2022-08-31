package iam

import (
	"fmt"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
)

const (
	userAPIVersion = "2"
)

// GetUserOptions describes search criteria for looking up users
type GetUserOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"organizationID,omitempty"`
	Name           *string `url:"name,omitempty"`
	LoginID        *string `url:"loginId,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	PageSize       *string `url:"pageSize,omitempty"`
	PageNumber     *string `url:"pageNumber,omitempty"`
	UserID         *string `url:"userId,omitempty"`
	ProfileType    *string `url:"profileType,omitempty" enum:"membership|accountStatus|passwordStatus|consentedApps|all"`
}

// UsersService provides operations on IAM User resources
type UsersService struct {
	client *Client

	validate *validator.Validate
}

// Parameters holds parameters
type Parameters struct {
	ResourceType string  `json:"resourceType"`
	Parameter    []Param `json:"parameter"`
}

// ChangeLoginIDRequest
type ChangeLoginIDRequest struct {
	LoginID string `json:"loginId"`
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
	UserUUIDs   []string
	PageNumber  int
	PageSize    int
	HasNextPage bool
}

// CreateUser creates a new IAM user.
func (u *UsersService) CreateUser(person Person) (*User, *Response, error) {
	if err := u.validate.Struct(person); err != nil {
		return nil, nil, err
	}
	req, err := u.client.newRequest(IDM, "POST", "authorize/identity/User", &person, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", "3")

	var bundleResponse interface{}

	doFunc := u.client.do
	if person.ManagingOrganization == "" { // Self registration
		doFunc = u.client.doSigned
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode() == http.StatusCreated {
		// Brand-new user
		var id string
		count, err := fmt.Sscanf(resp.Header.Get("Location"), "/authorize/identity/User/%s", &id)
		if err != nil {
			return nil, resp, ErrCouldNoReadResourceAfterCreate
		}
		if count == 0 {
			return nil, resp, ErrCouldNoReadResourceAfterCreate
		}
		return u.GetUserByID(id)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, resp, fmt.Errorf("unexpected StatusCode '%d' during user create", resp.StatusCode())
	}
	// HTTP 200
	return u.GetUserByID(person.LoginID)
}

// DeleteUser deletes the  IAM user.
func (u *UsersService) DeleteUser(person Person) (bool, *Response, error) {
	req, err := u.client.newRequest(IDM, "DELETE", "authorize/identity/User/"+person.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse interface{}

	doFunc := u.client.doSigned
	if !u.client.validSigner() {
		doFunc = u.client.do
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode() == http.StatusNoContent)
	return ok, resp, err
}

// RecoverPassword triggers the recovery flow for the given user
//
// Deprecated: Support end date is 1 Augustus 2020
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
	return u.userActionV(body, "$recover-password", "1")
}

// ChangeLoginID changes the loginID
// Link: https://www.hsdp.io/documentation/identity-and-access-management-iam/api-documents/resource-reference-api/user-api-v2#/User%20Management/post_User__id___change_loginid
func (u *UsersService) ChangeLoginID(user Person, newLoginID string) (bool, *Response, error) {
	body := &ChangeLoginIDRequest{
		LoginID: newLoginID,
	}
	req, err := u.client.newRequest(IDM, "POST", "authorize/identity/User/"+user.ID+"/$change-loginid", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}
	doFunc := u.client.doSigned
	if !u.client.validSigner() {
		doFunc = u.client.do
	}
	resp, err := doFunc(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode() == http.StatusNoContent)
	return ok, resp, err
}

// ResendActivation re-sends an activation email to the given user
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
	return u.userActionV(body, "$resend-activation", "2")
}

func (u *UsersService) userActionV(body *Parameters, action, apiVersion string) (bool, *Response, error) {
	req, err := u.client.newRequest(IDM, "POST", "authorize/identity/User/"+action, body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", apiVersion)

	var bundleResponse interface{}

	resp, err := u.client.doSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusAccepted)
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
	return u.userActionV(body, "$set-password", "2")
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
	return u.userActionV(body, "$change-password", "1")
}

func stringInc(number string) string {
	i, _ := strconv.Atoi(number)
	i = i + 1
	return strconv.FormatInt(int64(i), 10)
}

// GetAllUsers retrieves all users based on GetUserOptions
func (u *UsersService) GetAllUsers(opts *GetUserOptions, options ...OptionFunc) ([]string, *Response, error) {
	var users []string
	currentPage := "1"
	pageSize := "100"
	if opts.PageNumber == nil {
		opts.PageNumber = &currentPage
	}
	if opts.PageSize == nil {
		opts.PageSize = &pageSize
	}
	for {
		userList, resp, err := u.GetUsers(opts, options...)
		if err != nil {
			return users, resp, err
		}
		users = append(users, userList.UserUUIDs...)
		if !userList.HasNextPage {
			return users, resp, nil
		}
		// Next page
		currentPage = stringInc(currentPage)
		opts.PageNumber = &currentPage
	}
}

// GetUsers looks up users by search criteria specified in GetUserOptions
func (u *UsersService) GetUsers(opts *GetUserOptions, options ...OptionFunc) (*UserList, *Response, error) {
	req, err := u.client.newRequest(IDM, "GET", "security/users", opts, options)
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

	doFunc := u.client.doSigned
	if !u.client.validSigner() {
		doFunc = u.client.do
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return nil, resp, err
	}
	var list UserList

	list.HasNextPage = bundleResponse.Exchange.NextPageExists
	if opts != nil && opts.PageNumber != nil {
		pageNumber, _ := strconv.ParseInt(*opts.PageNumber, 10, 8)
		list.PageNumber = int(pageNumber)
	}
	for _, u := range bundleResponse.Exchange.Users {
		list.UserUUIDs = append(list.UserUUIDs, u.UserUUID)
	}

	return &list, resp, err
}

// GetUserByID looks up a user by UUID
func (u *UsersService) GetUserByID(uuid string) (*User, *Response, error) {
	opt := &GetUserOptions{
		UserID:      &uuid,
		ProfileType: String("all"),
	}
	req, _ := u.client.newRequest(IDM, "GET", "authorize/identity/User", opt, nil)
	req.Header.Set("api-version", "3")

	var responseStruct struct {
		Total int    `json:"total"`
		Entry []User `json:"entry"`
	}

	resp, err := u.client.do(req, &responseStruct)
	if err != nil {
		return nil, resp, fmt.Errorf("GetUserByID('%s'): %w", uuid, err)
	}
	if responseStruct.Total == 0 {
		return nil, resp, fmt.Errorf("GetUserByID('%s'): %w", uuid, ErrEmptyResults)
	}
	return &responseStruct.Entry[0], resp, nil
}

// GetUserIDByLoginID looks up the UUID of a user by LoginID (email address)
func (u *UsersService) GetUserIDByLoginID(loginID string) (string, *Response, error) {
	user, resp, err := u.GetUserByID(loginID)
	if err != nil {
		return "", resp, err
	}
	if user == nil {
		return "", resp, fmt.Errorf("GetUserIDByLoginID: %w", ErrEmptyResults)
	}
	return user.ID, resp, nil
}

// LegacyUpdateUser updates the user profile
func (u *UsersService) LegacyUpdateUser(profile Profile) (*Profile, *Response, error) {
	// don't send blank addresses
	profile.PruneBlankAddresses()
	// Also clear out un-settable fields
	profile.VerifiedMobilePhoneStatus = ""
	profile.EmailVerifiedStatus = ""
	profile.MustChangePassword = ""

	req, _ := u.client.newRequest(IDM, "PUT", "security/users/"+profile.ID, profile, nil)
	req.Header.Set("api-version", "2")

	var responseStruct struct {
		Exchange struct {
			UserUUID string  `json:"userUUID"`
			LoginID  string  `json:"loginId"`
			Profile  Profile `json:"profile"`
		} `json:"exchange"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}
	resp, err := u.client.do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	if responseStruct.ResponseCode != "200" {
		return nil, resp, fmt.Errorf("ResponseCode=%s: %v", responseStruct.ResponseCode, responseStruct)
	}
	return &responseStruct.Exchange.Profile, resp, err
}

// LegacyGetUserByUUID looks the a user by UUID using the legacy API
func (u *UsersService) LegacyGetUserByUUID(uuid string) (*Profile, *Response, error) {
	req, _ := u.client.newRequest(IDM, "GET", "security/users/"+uuid, nil, nil)
	req.Header.Set("api-version", userAPIVersion)

	var responseStruct struct {
		Exchange struct {
			LoginID string  `json:"loginId"`
			Profile Profile `json:"profile"`
		} `json:"exchange"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}
	resp, err := u.client.do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Exchange.Profile, resp, err
}

// LegacyGetUserIDByLoginID looks up the UUID of a user by LoginID (email address)
func (u *UsersService) LegacyGetUserIDByLoginID(loginID string) (string, *Response, error) {
	opt := &GetUserOptions{
		LoginID: &loginID,
	}
	req, _ := u.client.newRequest(IDM, "GET", "security/users", opt, nil)
	req.Header.Set("api-version", userAPIVersion)

	var responseStruct struct {
		Exchange struct {
			Users []struct {
				UserUUID string `json:"userUUID"`
			} `json:"users"`
		} `json:"exchange"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}

	resp, err := u.client.do(req, &responseStruct)
	if err != nil {
		return "", resp, err
	}
	if len(responseStruct.Exchange.Users) == 0 {
		return "", resp, fmt.Errorf("LegacyGetUserIDByLoginID: %w", ErrEmptyResults)
	}
	return responseStruct.Exchange.Users[0].UserUUID, resp, nil
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
	req, err := u.client.newRequest(IDM, "POST", "authorize/identity/User/"+userID+"/$mfa", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode() == http.StatusAccepted)
	return ok, resp, nil
}

// Unlock unlocks a user account with the given UserID
func (u *UsersService) Unlock(userID string) (bool, *Response, error) {
	req, err := u.client.newRequest(IDM, "POST", "authorize/identity/User/"+userID+"/$unlock", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode() == http.StatusNoContent)
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
