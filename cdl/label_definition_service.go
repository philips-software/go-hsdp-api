package cdl

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"path"
)

type LabelDefinition struct {
	ID           string            `json:"id,omitempty"`
	LabelDefName string            `json:"labelDefName" validate:"labelDefValidationHanlder"`
	Description  string            `json:"description,omitempty"`
	LabelScope   LabelScope        `json:"labelScope" validate:"labelDefValidationHanlder"`
	Label        string            `json:"label" validate:"required"`
	Type         string            `json:"type" validate:"required"`
	Labels       []LabelsArrayElem `json:"labels" validate:"labelDefValidationHanlder"`
	CreatedBy    string            `json:"createdBy,omitempty"`
	CreatedOn    string            `json:"createdOn,omitempty"`
}

type BundleEntry struct {
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource LabelDefinition `json:"resource,omitempty"`
}

type LabelDefBundleResponse struct {
	ResourceType string          `json:"resourceType,omitempty"`
	Id           string          `json:"id,omitempty"`
	Type         string          `json:"type,omitempty"`
	Link         json.RawMessage `json:"link,omitempty"`
	Entry        []BundleEntry   `json:"entry,required"`
}

type LabelScope struct {
	Type string `json:"type" validate:"required"`
}

type LabelsArrayElem struct {
	Label string `json:"label" validate:"required"`
}

type LabelDefinitionService struct {
	client   *Client
	config   *Config
	validate *validator.Validate
}

func (l *LabelDefinitionService) path(components ...string) string {
	return fmt.Sprintf("%s", path.Join(components...))
}

func (l *LabelDefinitionService) CreateLabelDefinition(studyId string, labelDef LabelDefinition) (*LabelDefinition, *Response, error) {
	if err := l.validate.Struct(labelDef); err != nil {
		return nil, nil, err
	}

	req, err := l.client.newCDLRequest("POST", l.path("Study", studyId, "LabelDef"), labelDef, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "1")

	var createdLabelDefinition LabelDefinition
	resp, err := l.client.do(req, &createdLabelDefinition)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateLabelDefinition: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdLabelDefinition, resp, nil
}

func (l *LabelDefinitionService) GetLabelDefinitions(studyId string, opt *GetOptions, options ...OptionFunc) ([]LabelDefinition, *Response, error) {
	req, err := l.client.newCDLRequest("GET", l.path("Study", studyId, "LabelDef"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "1")

	var getAllLabelDefResponse LabelDefBundleResponse
	resp, err := l.client.do(req, &getAllLabelDefResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	var labelDefinitionSlice []LabelDefinition
	for _, entry := range getAllLabelDefResponse.Entry {
		labelDefinitionSlice = append(labelDefinitionSlice, entry.Resource)
	}
	return labelDefinitionSlice, resp, err
}

func (l *LabelDefinitionService) GetLabelDefinitionByID(studyId string, labelDefId string) (*LabelDefinition, *Response, error) {
	req, err := l.client.newCDLRequest("GET", l.path("Study", studyId, "LabelDef", labelDefId), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "1")

	var labelDefinition LabelDefinition
	resp, err := l.client.do(req, &labelDefinition)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetLabelDefinitionByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &labelDefinition, resp, nil
}

func (l *LabelDefinitionService) DeleteLabelDefinitionById(studyId string, labelDefId string) (*Response, error) {
	req, err := l.client.newCDLRequest("DELETE", l.path("Study", studyId, "LabelDef", labelDefId), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", "1")

	resp, err := l.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("deleteLabelDefinitionById: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}
