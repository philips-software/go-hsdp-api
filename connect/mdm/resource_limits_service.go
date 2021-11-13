package mdm

import (
	"net/http"

	"github.com/philips-software/go-hsdp-api/internal"
)

type ResourceLimitsService struct {
	*Client
}

type ResourcesLimits map[string]int

func (r *ResourceLimitsService) get(which string) (*ResourcesLimits, *Response, error) {
	req, err := r.NewRequest(http.MethodGet, "/ResourcesLimit/"+which, nil)
	if err != nil {
		return nil, nil, err
	}
	var limits ResourcesLimits

	resp, err := r.Do(req, &limits)
	if err != nil {
		return nil, resp, err
	}
	if err := internal.CheckResponse(resp.Response); err != nil {
		return nil, resp, err
	}

	return &limits, resp, nil
}

func (r *ResourceLimitsService) GetDefault() (*ResourcesLimits, *Response, error) {
	return r.get("$default")
}

func (r *ResourceLimitsService) GetOverride() (*ResourcesLimits, *Response, error) {
	return r.get("$override")
}
