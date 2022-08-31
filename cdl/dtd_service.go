package cdl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/go-playground/validator/v10"
)

type DataTypeDefinition struct {
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	JsonSchema  json.RawMessage `json:"jsonSchema,omitempty"`
	CreatedOn   string          `json:"createdOn,omitempty"`
	CreatedBy   string          `json:"createdBy,omitempty"`
	UpdatedOn   string          `json:"updatedOn,omitempty"`
	UpdatedBy   string          `json:"updatedBy,omitempty"`
}

type DatatypeDefinitionService struct {
	client   *Client
	config   *Config
	validate *validator.Validate
}

func (dtd *DatatypeDefinitionService) path(components ...string) string {
	return path.Join(components...)
}

func (dtd *DatatypeDefinitionService) GetDataTypeDefinitions(opt *GetOptions, options ...OptionFunc) ([]DataTypeDefinition, *Response, error) {
	req, err := dtd.client.newCDLRequest("GET", dtd.path("DataTypeDefinition"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "3")

	var getAllDtdResponse []DataTypeDefinition
	resp, err := dtd.client.do(req, &getAllDtdResponse)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	return getAllDtdResponse, resp, err
}

func (dtd *DatatypeDefinitionService) CreateDataTypeDefinition(dataTypeDefinition DataTypeDefinition) (*DataTypeDefinition, *Response, error) {
	if err := dtd.validate.Struct(dataTypeDefinition); err != nil {
		return nil, nil, err
	}
	req, err := dtd.client.newCDLRequest("POST", dtd.path("DataTypeDefinition"), dataTypeDefinition, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "3")

	var createdDtd DataTypeDefinition
	resp, err := dtd.client.do(req, &createdDtd)

	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateDataTypeDefinition: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdDtd, resp, nil
}

func (dtd *DatatypeDefinitionService) GetDataTypeDefinitionByID(id string) (*DataTypeDefinition, *Response, error) {
	req, err := dtd.client.newCDLRequest("GET", dtd.path("DataTypeDefinition", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "3")

	var dtdByIdResponse DataTypeDefinition
	resp, err := dtd.client.do(req, &dtdByIdResponse)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	return &dtdByIdResponse, resp, err
}

func (dtd *DatatypeDefinitionService) UpdateDataTypeDefinition(dataTypeDefinition DataTypeDefinition) (*DataTypeDefinition, *Response, error) {
	req, err := dtd.client.newCDLRequest("PUT", dtd.path("DataTypeDefinition", dataTypeDefinition.ID),
		dataTypeDefinition, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "3")

	var updatedDtd DataTypeDefinition
	resp, err := dtd.client.do(req, &updatedDtd)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("UpdateDataTypeDefinition: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &updatedDtd, resp, nil
}
