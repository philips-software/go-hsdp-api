package cdl_test

import (
	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestExportRouteCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	exportRouteID := "1ee7eb94-3eab-4c1c-9a1c-7c5347ab538d"

	exportRoute := cdl.ExportRoute{
		ExportRouteName: "ExportTrial_for_demo27",
		Description:     "description11",
		DisplayName:     "DisplayName",
		Source: cdl.Source{
			CDLResearchStudy: cdl.ExportResearchStudySource{
				Endpoint: "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa",
			},
		},
		Destination: cdl.Destination{
			CDLResearchStudy: cdl.ExportResearchStudyDestination{
				Endpoint: "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/5c8431e2-f4f1-11eb-bf8f-b799651c8a11",
			},
		},
		ServiceAccount: cdl.ExportServiceAccount{
			CDLServiceAccount: cdl.ExportServiceAccountDetails{
				ServiceID:           "eng_cdl_tenant_1_ser.eng__cdl__tenant__1__app.eng__cdl__tenant__1@eng__cdl__tenant__1.cdal.philips-healthsuite.com",
				PrivateKey:          "-----BEGIN RSA PRIVATE KEY-----SERVICE_KEY-----END RSA PRIVATE KEY-----",
				AccessTokenEndPoint: "https://iam-development.us-east.philips-healthsuite.com/oauth2/access_token",
				TokenEndPoint:       "https://iam-development.us-east.philips-healthsuite.com/authorize/oauth2/token",
			},
		},
	}

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/ExportRoute", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
      "id": "1ee7eb94-3eab-4c1c-9a1c-7c5347ab538d",
      "resourceType": "ExportRoute",
      "name": "ExportTrial_for_demo74",
      "description": "description11",
      "displayName": "display name11",
      "source": {
        "cdlResearchStudy": {
          "endpoint": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa"
        }
      },
      "destination": {
        "cdlResearchStudy": {
          "endpoint": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/5c8431e2-f4f1-11eb-bf8f-b799651c8a11"
        }
      },
      "autoExport": false,
      "createdBy": "user@philips.com",
      "createdOn": "2021-08-04T10:19:53.168+00:00",
      "updatedBy": "user@philips.com",
      "updatedOn": "2021-08-04T10:19:53.168+00:00"
    }
  ]
}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "id": "b6d0c8e6-6df4-42ae-a239-898e92c5814d",
  "type": "searchset",
  "link": [
    {
      "relation": "self",
      "url": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa/LabelDef"
    }
  ],
  "entry": [
    {
      "fullUrl": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/ExportRoute/1ee7eb94-3eab-4c1c-9a1c-7c5347ab538d",
      "resource": {
        "id": "1ee7eb94-3eab-4c1c-9a1c-7c5347ab538d",
        "resourceType": "ExportRoute",
        "name": "ExportTrial_for_demo74",
        "description": "description11",
        "displayName": "display name11",
        "source": {
          "cdlResearchStudy": {
            "endpoint": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa"
          }
        },
        "destination": {
          "cdlResearchStudy": {
            "endpoint": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/5c8431e2-f4f1-11eb-bf8f-b799651c8a11"
          }
        },
        "autoExport": false,
        "createdBy": "user@philips.com",
        "createdOn": "2021-08-04T10:19:53.168+00:00",
        "updatedBy": "user@philips.com",
        "updatedOn": "2021-08-04T10:19:53.168+00:00"
      }
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/ExportRoute/"+exportRouteID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	createdExportRoute, resp, err := cdlClient.ExportRoute.CreateExportRoute(exportRoute)

	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, createdExportRoute) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, createdExportRoute.ID, exportRouteID)

	item, resp, err := cdlClient.ExportRoute.GetExportRouteByID(exportRouteID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, exportRouteID, item.ID)

	items, bundleResponse, resp, err := cdlClient.ExportRoute.GetExportRoutes(1, nil)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, bundleResponse) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, exportRouteID, items[0].ID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = cdlClient.ExportRoute.DeleteExportRouteByID(exportRouteID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
