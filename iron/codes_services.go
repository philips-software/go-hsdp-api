package iron

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type CodesServices struct {
	client    *Client
	token     string
	projectID string
}

// Code describes a Iron code package
type Code struct {
	ID              string     `json:"id,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	ProjectID       string     `json:"project_id,omitempty"`
	Name            string     `json:"name"`
	Image           string     `json:"image"`
	LatestChecksum  string     `json:"latest_checksum,omitempty"`
	Rev             int        `json:"rev,omitempty"`
	LatestHistoryID string     `json:"latest_history_id,omitempty"`
	LatestChange    *time.Time `json:"latest_change,omitempty"`
}

// DockerCredentials describes a set of docker credentials
type DockerCredentials struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ServerAddress string `json:"serveraddress"`
}

func (d DockerCredentials) Valid() bool {
	if d.Email == "" || d.Username == "" || d.Password == "" || d.ServerAddress == "" {
		return false
	}
	return true
}

// CreateOrUpdateCode creates or updates code packages on Iron which can be used to run tasks
func (c *CodesServices) CreateOrUpdateCode(code Code) (*Code, *Response, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	j, _ := json.Marshal(code)
	r := bytes.NewReader(j)
	w := multipart.NewWriter(&b)
	if fw, err = w.CreateFormField("data"); err != nil {
		return nil, nil, err
	}
	if _, err = io.Copy(fw, r); err != nil {
		return nil, nil, err
	}
	_ = w.Close()

	req, err := http.NewRequest("POST", c.client.baseIRONURL.String()+"/2/projects/"+c.projectID+"/codes", &b)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "OAuth "+c.token)

	var createResponse struct {
		Message string `json:"msg"`
		ID      string `json:"id,omitempty"`
	}
	resp, err := c.client.do(req, &createResponse)
	if err != nil {
		return nil, resp, err
	}
	if createResponse.ID == "" {
		return nil, resp, fmt.Errorf("empty code value: '%s'", createResponse.Message)
	}
	return c.GetCode(createResponse.ID)
}

func (c *CodesServices) GetCodes() (*[]Code, *Response, error) {
	getPath := c.client.Path("projects", c.projectID, "codes")

	page := 0
	perPage := 100

	req, err := c.client.newRequest(
		"GET",
		getPath,
		pageOptions{
			Page:    &page,
			PerPage: &perPage,
		},
		nil)
	if err != nil {
		return nil, nil, err
	}
	var codes struct {
		Codes []Code `json:"codes"`
	}
	resp, err := c.client.do(req, &codes)
	return &codes.Codes, resp, err
}

func (c *CodesServices) GetCode(codeID string) (*Code, *Response, error) {
	req, err := c.client.newRequest(
		"GET",
		c.client.Path("projects", c.projectID, "codes", codeID),
		nil,
		nil)
	if err != nil {
		return nil, nil, err
	}
	var code Code
	resp, err := c.client.do(req, &code)
	return &code, resp, err
}

// DeleteCode deletes a code from Iron
func (c *CodesServices) DeleteCode(codeID string) (bool, *Response, error) {
	req, err := c.client.newRequest(
		"DELETE",
		c.client.Path("projects", c.projectID, "codes", codeID),
		nil,
		nil)
	if err != nil {
		return false, nil, err
	}
	var deleteResponse struct {
		Message string `json:"msg,omitempty"`
	}
	resp, err := c.client.do(req, &deleteResponse)
	if deleteResponse.Message != "Deleted" {
		return false, resp, err
	}
	return true, resp, nil
}

// DockerLogin stores private Docker registry credentials so Iron can fetch images when needed
func (c *CodesServices) DockerLogin(creds DockerCredentials) (bool, *Response, error) {
	if !creds.Valid() {
		return false, nil, ErrInvalidDockerCredentials
	}
	data, err := json.Marshal(&creds)
	if err != nil {
		return false, nil, err
	}
	authString := base64.StdEncoding.EncodeToString(data)
	var authRequest struct {
		Auth string `json:"auth"`
	}
	authRequest.Auth = authString
	req, err := c.client.newRequest(
		"POST",
		c.client.Path("projects", c.projectID, "credentials"),
		&authRequest,
		nil)
	if err != nil {
		return false, nil, err
	}
	var authResponse struct {
		Message string `json:"msg"`
	}
	var success bool
	resp, err := c.client.do(req, &authResponse)
	if resp != nil {
		success = resp.StatusCode == http.StatusOK
	}
	return success, resp, err
}
