package has

import (
	"net/http"
)

// SessionsService provides operations on HAS sessions
type SessionsService struct {
	orgID  string
	client *Client
}

// Session describes a HAS session
type Session struct {
	// reference to session (to query specific status)
	SessionID string `json:"sessionId,omitempty"`
	//Full URL to resource (empty if not yet available)
	SessionURL string `json:"sessionUrl,omitempty"`
	// The session type
	SessionType string `json:"sessionType,omitempty" enum:"DEV"`
	// Enumerated status of resource claim
	State string `json:"state,omitempty" enum:"PENDING|AVAILABLE|TIMEDOUT|INUSE"`
	// The id of the resource that is used for this session
	ResourceID string `json:"resourceId,omitempty"`
	// The region where the resource was provisioned.
	Region string `json:"region,omitempty"`
	// The id of the user that has claimed this session.
	UserID string `json:"userId,omitempty"`
	// The image id of the session
	ImageID string `json:"imageId,omitempty"`
	// The cluster tag to target the specific cluster
	ClusterTag string `json:"clusterTag,omitempty"`
	// The remote IP of the instance
	RemoteIP string `json:"remoteIp,omitempty"`
	// The access token for the instance
	AccessToken string `json:"accessToken,omitempty""`
}

// Sessions contains a list of Session values
type Sessions struct {
	Sessions []Session `json:"sessions"`
	ResponseError
}

// ResponseError describes the response fields in case of error
type ResponseError struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// CreateSession creates a new user session in HAS
func (c *SessionsService) CreateSession(userID string, session Session) (*Sessions, *Response, error) {
	req, err := c.client.newHASRequest("POST", "user/"+userID+"/session", &session, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", APIVersion)

	var sr Sessions
	resp, err := c.client.do(req, &sr)
	if err != nil {
		return nil, resp, err
	}
	return &sr, resp, nil
}

// SessionOptions describes options (query) parameters which
// can be used on some API endpoints
type SessionOptions struct {
	ResourceID *string `url:"resourceId,omitempty"`
}

// GetSession gets a user session in HAS
func (c *SessionsService) GetSession(userID string, opt *SessionOptions) (*Sessions, *Response, error) {
	req, err := c.client.newHASRequest("GET", "user/"+userID+"/session", &opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", APIVersion)

	var sr Sessions
	resp, err := c.client.do(req, &sr)
	if err != nil {
		return nil, resp, err
	}
	return &sr, resp, nil
}

// GetSessions gets all sessions in HAS
func (c *SessionsService) GetSessions() (*Sessions, *Response, error) {
	req, err := c.client.newHASRequest("GET", "session", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", APIVersion)

	var sr Sessions
	resp, err := c.client.do(req, &sr)
	if err != nil {
		return nil, resp, err
	}
	return &sr, resp, nil
}

// DeleteSession deletes a user session in HAS
func (c *SessionsService) DeleteSession(userID string) (bool, *Response, error) {
	req, err := c.client.newHASRequest("DELETE", "user/"+userID+"/session", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", APIVersion)

	var sr Sessions
	resp, _ := c.client.do(req, &sr)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, nil, ErrEmptyResults
	}
	return true, resp, nil
}
