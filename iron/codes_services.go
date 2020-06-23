package iron

import (
	"bytes"
	"encoding/json"
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
	LatestChecksum  string     `json:"latest_checksum"`
	Rev             int        `json:"rev,omitempty"`
	LatestHistoryID string     `json:"latest_history_id,omitempty"`
	LatestChange    time.Time  `json:"latest_change,omitempty"`
}

func (c *CodesServices) CreateOrUpdateCode(code Code) (*Code, *Response, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	j, err := json.Marshal(code)
	if err != nil {
		return nil, nil, err
	}
	r := bytes.NewReader(j)
	w := multipart.NewWriter(&b)
	if fw, err = w.CreateFormField("data"); err != nil {
		return nil, nil, err
	}
	if _, err = io.Copy(fw, r); err != nil {
		return nil, nil, err
	}
	_ = w.Close()

	req, err := http.NewRequest("POST", c.client.baseIRONURL.String()+"projects/"+c.projectID+"/codes", &b)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "OAuth "+c.token)

	var createResponse struct {
		Message string `json:"msg,omitempty"`
		ID      string `json:"id,omitempty"`
	}
	resp, err := c.client.Do(req, &createResponse)
	if err != nil {
		return nil, resp, err
	}
	return c.GetCode(createResponse.ID)
}

func (c *CodesServices) GetCodes() (*[]Code, *Response, error) {
	req, err := c.client.NewRequest("GET", "projects/"+c.projectID+"/codes", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var codes struct {
		Codes []Code `json:"codes"`
	}
	resp, err := c.client.Do(req, &codes)
	return &codes.Codes, resp, err
}

func (c *CodesServices) GetCode(codeID string) (*Code, *Response, error) {
	req, err := c.client.NewRequest("GET", "projects/"+c.projectID+"/codes/"+codeID, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var code Code
	resp, err := c.client.Do(req, &code)
	return &code, resp, err
}

func (c *CodesServices) DeleteCode(codeID string) (bool, *Response, error) {
	req, err := c.client.NewRequest("DELETE", "projects/"+c.projectID+"/codes/"+codeID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	var deleteResponse struct {
		Message string `json:"msg,omitempty"`
	}
	resp, err := c.client.Do(req, &deleteResponse)
	if deleteResponse.Message != "Deleted" {
		return false, resp, err
	}
	return true, resp, nil
}
