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
			CdlResearchStudy: cdl.ExportResearchStudySource{
				Endpoint: "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa",
			},
		},
		Destination: cdl.Destination{
			CdlResearchStudy: cdl.ExportResearchStudyDestination{
				Endpoint: "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/5c8431e2-f4f1-11eb-bf8f-b799651c8a11",
			},
		},
		ServiceAccount: cdl.ExportServiceAccount{
			CdlServiceAccount: cdl.ExportServiceAccountDetails{
				ServiceId:           "eng_cdl_tenant_1_ser.eng__cdl__tenant__1__app.eng__cdl__tenant__1@eng__cdl__tenant__1.cdal.philips-healthsuite.com",
				PrivateKey:          "-----BEGIN RSA PRIVATE KEY-----MIIEogIBAAKCAQEAiHEAGs0T2aadjxXVjaBjRMoAmPBMWgFXdGSfP4FNQdNxIZKcSqPcnTihLQdQL/ZZh4mlG0iGecQLCNkEGDN2/rFmDBUfpObR6Bbhfe9IW1ikEzADOOvVaZiKzOBGBSfIV7hUtpWFLeY8ehyWdriF84e05hq+zeG24zxpP1nVg8iAvqwy8tZG+tokvXhTNREY80laf2f7rou1JVuN5awj6L1jnxllnllwQCqi0OtPL029RAd9o20qdTcsazebf3HMTA0KZE7fNhoyNpExXwA5a5DHAL0IwRuMjzZmOyNvtE+RI9FcmTyTmxxXddtU5/MkE7q3TVz2wtyZ3Wwbzi0ZeQIDAQABAoIBAGILgZ3AvIDsv8/wSrMpC+yElAVSuCr9R9D19ZF24wNoY6VSa0kfkwrPhNKMrMyqZA+Hk8OVTDz36U4EVHLnmJzQ8ZnBHCotS61Rs9wBNKeffNfk6ovevE7TNPqgNvDBc6/FX+qMa1VeWxsMq/rIAknLvOyVT4M2rWuBH8hMT6gQPUTzesQ2LGon51AawS+/UXzSOHwqvltH/zphPv4jhtCAkY7qKvLqogOjPU4d/UefS0ZK7/4TmAdGeJlMiv26B9/getBEd7Qk3W+f2YwiLWtKs2d2YkAqlTqpyKE0KwChw3/Whuj3g1CINPVOfFfliv/MKLmjSBFTeJyzGcN823kCgYEA4ZWVysi+NvxQjFJNKGTjZeYWmRySezS7nd+G6YACoIu6c0O4xC3TQsfhMfVmjAM1ffu7eh1I2MFGc+Yo6qBNtqYjT7Hlm9a6feijoM94MSJUDkQon6B5BvhjSRMlXaR6fZb0w4XjYaUNDJVes79RY1nnAR1VoauRqZCnnOhIWpcCgYEAmtaALk5V8Rg52mOcv+K2WibIB652gu7ICv8S62/xsntbtRjgLiBmlDJc8JsCPOrnYlXoSNZpQdNCbHCL/NfDpzkVaKHH990zqiW9Edwln2LiH13Bo6MMn6zE+/sTkCkNCH7m7tt/T3eZx/tw9H1TCR1Jl7Fvv1lwkQoiqlP0/m8CgYAxg7TiUte1mAJSGoqHEEX9ithw+R2J35RC3dpuDEQHW0QsorO+k9RoNxlN7vB4UQf/xC5talogAaRmMiHPBiqoqaTcjE66uxIqKtMnrAJUvpU2oG1ORFnsVr7sPkCYYk7knCrTc+Lp/uFzXqHv0FGb/hK/YuH134PUdUTlIvMmtwKBgEuvAWShAb0hHFY1To80n/Gc9zVZ/6+sS7ekSnkudLPLPF5e1GV3jOxvWaJ6AjQIliUo3KuNslFslBExShvC023PpzlHqtjrFK/cVnh+ZR1tVh4C0/3KWwdJidepODzE9AvtC7BBNg9/5Hkt3F6FS6su16QAJSEg9LbQf3VGKICdAoGARoL7lDFrijgz/98fKYItRNafuui+sgtb5dAFGFYMDcgyowATVxk4AMvocPhyPA3SSPIcvJs1qlIhAzwI6IG8mNlrtIsZo0P+7ybtRWO+ATg+pmOJJ+00ph0WY78ptAraUS7vcF+0XaibLjetYdNMJEcn6kBtoUR20ZUYSLwG1Gc=-----END RSA PRIVATE KEY-----",
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
			io.WriteString(w, `{
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
			io.WriteString(w, `{
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

	resp, err = cdlClient.ExportRoute.DeleteExportRouteById(exportRouteID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
