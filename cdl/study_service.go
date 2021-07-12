package cdl

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

type Period struct {
	End string `json:"end"`
}

type Study struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	StudyOwner  string `json:"studyOwner" validate:"required"`
	Period      Period `json:"period" validate:"required"`
}

type StudyService struct {
	client   *Client
	config   *Config
	validate *validator.Validate
}

func (s *StudyService) path(remainder string) string {
	return fmt.Sprintf("store/cdl/%s/%s", s.config.OrganizationID, remainder)
}

func (s *StudyService) CreateStudy(study Study) (*Study, *Response, error) {
	if err := s.validate.Struct(study); err != nil {
		return nil, nil, err
	}
	req, err := s.client.newCDLRequest("POST", s.path("Study"), study, nil)
	if err != nil {
		return nil, nil, err
	}
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
