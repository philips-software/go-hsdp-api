package inference

import (
	"bytes"
	"fmt"
	"io"
	"path"

	"github.com/go-playground/validator/v10"
)

type ComputeProviderService struct {
	client *Client

	validate *validator.Validate
}

type UpdateRequest struct {
	AccessKey string `json:"accessKey" validate:"required"`
	SecretKey string `json:"secretKey" validate:"required"`
}

func (s *ComputeProviderService) path(components ...string) string {
	return path.Join(components...)
}

func (s *ComputeProviderService) UpdateProvider(request UpdateRequest) (bool, *Response, error) {
	if err := s.validate.Struct(request); err != nil {
		return false, nil, err
	}
	req, err := s.client.newInferenceRequest("POST", s.path("ComputeProvider"), request, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Api-Version", APIVersion)

	var operationOutcome bytes.Buffer
	resp, err := s.client.do(req, &operationOutcome)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("UpdateProvider: %w", ErrEmptyResult)
		}
		return false, resp, err
	}
	return true, resp, nil
}
