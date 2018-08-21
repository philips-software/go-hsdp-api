package iam

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/jeffail/gabs"
)

const servicesAPIVersion = "1"

// Service represents a IAM service resource
type Service struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	ServiceID      string `json:"serviceId,omitempty"`
	Description    string `json:"description"`
	ApplicationID  string `json:"applicationId"`
	OrganizationID string `json:"organizationId,omitempty"`
	ExpiresOn      string `json:"expiresOn,omitempty"`
	PrivateKey     string `json:"privateKey,omitempty"`
}

// ServicesService provides operations on IAM Sessions resources
type ServicesService struct {
	client *Client
}

// GetServiceOptions describes search criteria for looking up services
type GetServiceOptions struct {
	ID             *string `url:"_id,omitempty"`
	Name           *string `url:"name,omitempty"`
	ApplicationID  *string `url:"applicationId,omitempty"`
	OrganizationID *string `url:"organizationId,omitempty"`
}

// GetServiceByID looks up a service by ID
func (p *ServicesService) GetServiceByID(id string) (*Service, *Response, error) {
	return p.GetService(&GetServiceOptions{ID: &id}, nil)
}

// GetServiceByName looks up a service by name
func (p *ServicesService) GetServiceByName(name string) (*Service, *Response, error) {
	return p.GetService(&GetServiceOptions{Name: &name}, nil)
}

// GetServicesByApplicationID finds all services which belong to the applicationID
func (p *ServicesService) GetServicesByApplicationID(applicationID string) (*[]Service, *Response, error) {
	opt := &GetServiceOptions{
		ApplicationID: &applicationID,
	}
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Service", opt, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)

	var responseStruct struct {
		Total int       `json:"total"`
		Entry []Service `json:"entry"`
	}

	resp, err := p.client.Do(req, &responseStruct)
	if err != nil {
		return nil, resp, err
	}
	return &responseStruct.Entry, resp, err
}

// CreateService creates a Service
func (p *ServicesService) CreateService(name, description, applicationID string) (*Service, *Response, error) {
	role := &Service{
		Name:          name,
		Description:   description,
		ApplicationID: applicationID,
	}
	req, err := p.client.NewRequest(IDM, "POST", "authorize/identity/Service", role, nil)
	req.Header.Set("api-version", servicesAPIVersion)

	var createdService Service

	resp, err := p.client.Do(req, &createdService)
	if err != nil {
		return nil, resp, err
	}
	return &createdService, resp, err
}

// GetService looks up a services based on GetServiceOptions
func (p *ServicesService) GetService(opt *GetServiceOptions, options ...OptionFunc) (*Service, *Response, error) {
	services, resp, err := p.GetServices(opt, options...)
	if err != nil {
		return nil, resp, err
	}
	return &(*services)[0], resp, nil
}

// GetServices looks up services based on GetServiceOptions
func (p *ServicesService) GetServices(opt *GetServiceOptions, options ...OptionFunc) (*[]Service, *Response, error) {
	req, err := p.client.NewRequest(IDM, "GET", "authorize/identity/Service", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)

	var bundleResponse bytes.Buffer

	resp, err := p.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	services, err := p.parseFromBundle(bundleResponse.Bytes())
	return services, resp, err
}

// DeleteService deletes the given Service
func (p *ServicesService) DeleteService(service Service) (bool, *Response, error) {
	req, err := p.client.NewRequest(IDM, "DELETE", "authorize/identity/Service/"+service.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", servicesAPIVersion)

	var deleteResponse interface{}

	resp, err := p.client.Do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, nil
	}
	return true, resp, err
}

func (p *ServicesService) parseFromBundle(bundle []byte) (*[]Service, error) {
	jsonParsed, err := gabs.ParseJSON(bundle)
	if err != nil {
		return nil, err
	}
	count, ok := jsonParsed.S("total").Data().(float64)
	if !ok || count == 0 {
		return nil, errors.New("empty result")
	}
	services := make([]Service, int64(count))

	children, _ := jsonParsed.S("entry").Children()
	for i, r := range children {
		var s Service
		s.ID = r.Path("id").Data().(string)
		s.Name, _ = r.Path("name").Data().(string)
		s.Description, _ = r.Path("description").Data().(string)
		s.ServiceID, _ = r.Path("serviceId").Data().(string)
		s.OrganizationID, _ = r.Path("organizationId").Data().(string)
		s.ApplicationID, _ = r.Path("applicationId").Data().(string)
		s.PrivateKey, _ = r.Path("privateKey").Data().(string)
		s.ExpiresOn, _ = r.Path("expiresOn").Data().(string)
		services[i] = s
	}
	return &services, nil
}
