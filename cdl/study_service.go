package cdl

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/go-playground/validator/v10"
)

type Period struct {
	End string `json:"end"`
}

type Study struct {
	ID                        string `json:"id,omitempty"`
	Title                     string `json:"title" validate:"required"`
	Description               string `json:"description,omitempty"`
	Organization              string `json:"organization,omitempty"`
	StudyOwner                string `json:"studyOwner" validate:"required"`
	Period                    Period `json:"period" validate:"required"`
	DataProtectedFromDeletion bool   `json:"dataProtectedFromDeletion,omitempty"`
}

type StudyService struct {
	client   *Client
	config   *Config
	validate *validator.Validate
}

// GetOptions describes the fields on which you can search for studies
type GetOptions struct {
}

func (s *StudyService) path(components ...string) string {
	return fmt.Sprintf("/%s", path.Join(components...))
}

func (s *StudyService) CreateStudy(study Study) (*Study, *Response, error) {
	if err := s.validate.Struct(study); err != nil {
		return nil, nil, err
	}
	req, err := s.client.newCDLRequest("POST", s.path("Study"), study, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "2")

	var createdStudy Study
	resp, err := s.client.do(req, &createdStudy)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateStudy: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &createdStudy, resp, nil
}

func (s *StudyService) GetStudyByID(id string) (*Study, *Response, error) {
	req, err := s.client.newCDLRequest("GET", s.path("Study", id), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "2")

	var foundStudy Study
	resp, err := s.client.do(req, &foundStudy)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("GetStudyByID: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &foundStudy, resp, nil
}

func (s *StudyService) UpdateStudy(study Study) (*Study, *Response, error) {
	req, err := s.client.newCDLRequest("PUT", s.path("Study", study.ID), study, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "2")

	var updated Study
	resp, err := s.client.do(req, &updated)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("UpdateStudy: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}
	return &updated, resp, nil
}

func (s *StudyService) GetStudies(opt *GetOptions, options ...OptionFunc) ([]Study, *Response, error) {
	var studies []Study

	req, err := s.client.newCDLRequest("GET", s.path("Study"), opt, options...)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "3")

	var bundleResponse struct {
		ResourceType string  `json:"resourceType,omitempty"`
		Type         string  `json:"type,omitempty"`
		Total        int     `json:"total"`
		Entry        []Study `json:"entry"`
	}
	resp, err := s.client.do(req, &bundleResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, ErrEmptyResult
		}
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return studies, resp, ErrEmptyResult
	}
	for _, e := range bundleResponse.Entry {
		studies = append(studies, e)
	}
	return studies, resp, err
}
