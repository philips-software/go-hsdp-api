package cdl

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"path"
)

type ExportRoute struct {
	ID              string               `json:"id,omitempty"`
	ExportRouteName string               `json:"name" validate:"required"`
	Description     string               `json:"description,omitempty"`
	DisplayName     string               `json:"displayName" validate:"required"`
	Source          Source               `json:"source" validate:"required"`
	AutoExport      bool                 `json:"autoExport,omitempty"`
	Destination     Destination          `json:"destination" validate:"required"`
	ServiceAccount  ExportServiceAccount `json:"serviceAccount" validate:"required"`
	CreatedBy       string               `json:"createdBy,omitempty"`
	CreatedOn       string               `json:"createdOn,omitempty"`
	UpdatedBy       string               `json:"updatedBy,omitempty"`
	UpdatedOn       string               `json:"updatedOn,omitempty"`
}

type Source struct {
	CdlResearchStudy ExportResearchStudySource `json:"cdlResearchStudy" validate:"required"`
}

type Destination struct {
	CdlResearchStudy ExportResearchStudyDestination `json:"cdlResearchStudy" validate:"required"`
}

type ExportResearchStudySource struct {
	Endpoint string              `json:"endpoint" validate:"required"`
	Allowed  *ExportAllowedField `json:"allowed,omitempty"`
}

type ExportResearchStudyDestination struct {
	Endpoint string `json:"endpoint" validate:"required"`
}

type ExportAllowedField struct {
	DataObject []ExportDataObject `json:"dataObject,omitempty"`
}

type ExportDataObject struct {
	Type        string        `json:"type"`
	ExportLabel []ExportLabel `json:"label,omitempty"`
}

type ExportLabel struct {
	Name             string `json:"name,omitempty"`
	ApprovalRequired bool   `json:"approvalRequired,omitempty"`
}

type ExportServiceAccount struct {
	CdlServiceAccount ExportServiceAccountDetails `json:"cdlServiceAccount" validate:"required"`
}

type ExportServiceAccountDetails struct {
	ServiceId           string `json:"serviceId" validate:"required"`
	PrivateKey          string `json:"privateKey" validate:"required"`
	AccessTokenEndPoint string `json:"accessTokenEndPoint" validate:"required"`
	TokenEndPoint       string `json:"tokenEndPoint" validate:"required"`
}

type ExportRouteService struct {
	client   *Client
	config   *Config
	validate *validator.Validate
}

func (exp *ExportRouteService) path(components ...string) string {
	return path.Join(components...)
}

type ExportRouteBundleEntry struct {
	FullURL  string      `json:"fullUrl,omitempty"`
	Resource ExportRoute `json:"resource,omitempty"`
}

type LinkElementType struct {
	Relation string `json:"relation"`
	Url      string `json:"url"`
}

type ExportRouteBundleResponse struct {
	ResourceType string                   `json:"resourceType,omitempty"`
	Id           string                   `json:"id,omitempty"`
	Type         string                   `json:"type,omitempty"`
	Link         []LinkElementType        `json:"link,omitempty"`
	Entry        []ExportRouteBundleEntry `json:"entry,required"`
}

func (exp *ExportRouteService) CreateExportRoute(exportRoute ExportRoute) (*ExportRoute, *Response, error) {
	if err := exp.validate.Struct(exportRoute); err != nil {
		return nil, nil, err
	}

	req, err := exp.client.newCDLRequest("POST", exp.path("ExportRoute"), exportRoute, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", "1")

	var createdExportRoute ExportRoute
	resp, err := exp.client.do(req, &createdExportRoute)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("CreateExportRoute: %w", ErrEmptyResult)
		}
		return nil, resp, err
	}

	return &createdExportRoute, resp, nil
}

func (exp *ExportRouteService) GetExportRoutes(page int, options ...OptionFunc) ([]ExportRoute, *ExportRouteBundleResponse, *Response, error) {
	req, err := exp.client.newCDLRequest("GET", exp.path("ExportRoute"), &struct {
		Page int `url:"page"`
	}{page}, options...)
	if err != nil {
		return nil, nil, nil, err
	}
	req.Header.Set("Api-Version", "1")

	var getAllExportRouteResponse ExportRouteBundleResponse
	resp, err := exp.client.do(req, &getAllExportRouteResponse)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil, resp, ErrEmptyResult
		}
		return nil, nil, resp, err
	}
	var exportRouteSlice []ExportRoute
	for _, entry := range getAllExportRouteResponse.Entry {
		exportRouteSlice = append(exportRouteSlice, entry.Resource)
	}
	return exportRouteSlice, &getAllExportRouteResponse, resp, err
}

func (exp *ExportRouteService) GetExportRouteByID(exportRouteId string) (*ExportRoute, *Response, error) {
	page := 1
	exportRoutes, getAllExportBundleResponse, resp, err := exp.GetExportRoutes(page)

	if err != nil {
		return nil, resp, err
	}

	for {
		for _, expRoute := range exportRoutes {
			if expRoute.ID == exportRouteId {
				return &expRoute, resp, err
			}
		}

		lastPage := true
		for _, link := range getAllExportBundleResponse.Link {
			if link.Relation == "next" {
				page += 1
				exportRoutes, getAllExportBundleResponse, resp, err = exp.GetExportRoutes(page)
				if err != nil {
					return nil, resp, err
				}
				lastPage = false
			}
		}
		if lastPage {
			return nil, resp, err
		}
	}
}

func (exp *ExportRouteService) DeleteExportRouteById(exportRouteId string) (*Response, error) {
	req, err := exp.client.newCDLRequest("DELETE", exp.path("ExportRoute", exportRouteId), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Version", "1")

	resp, err := exp.client.do(req, nil)
	if (err != nil && err != io.EOF) || resp == nil {
		if resp == nil && err != nil {
			err = fmt.Errorf("DeleteExportRouteById: %w", ErrEmptyResult)
		}
		return resp, err
	}
	return resp, nil
}
