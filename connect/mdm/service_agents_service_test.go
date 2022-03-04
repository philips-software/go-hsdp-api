package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceAgentsGet(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxMDM.HandleFunc("/connect/mdm/ServiceAgent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "meta": {
    "lastUpdated": "2022-03-04T10:24:51.682009+00:00"
  },
  "id": "ab34b6f3-7735-4176-9e1f-6f34492d1b6e",
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "search": {
        "mode": "global"
      },
      "resource": {
        "meta": {
          "lastUpdated": "2018-11-27T12:02:46.728293+00:00",
          "versionId": "a2246795-c764-4f00-a3db-2213a72f8ce9"
        },
        "id": "a2246795-c764-4f00-a3db-2213a72f8ce9",
        "resourceType": "ServiceAgent",
        "name": "blobserviceagent",
        "description": "ServiceAgentBlobPublishing",
        "apiVersionSupported": "1",
        "dataSubscriberId": {
          "reference": "DataSubscriber/501d5f20-707c-4a6b-9adc-673be25ddd6b"
        },
        "authenticationMethodIds": [
          {
            "reference": "AuthenticationMethod/adb099e6-a7a3-41e2-913d-82b96c6fed9f"
          }
        ]
      },
      "fullUrl": "ServiceAgent/a2246795-c764-4f00-a3db-2213a72f8ce9"
    },
    {
      "search": {
        "mode": "global"
      },
      "resource": {
        "meta": {
          "lastUpdated": "2019-01-31T12:46:54.604356+00:00",
          "versionId": "8da12d84-4e4c-4c47-bdec-26b3614f6acf"
        },
        "id": "8da12d84-4e4c-4c47-bdec-26b3614f6acf",
        "resourceType": "ServiceAgent",
        "name": "postgreserviceagent",
        "description": "ServiceAgentPostgrePublishing",
        "apiVersionSupported": "1",
        "dataSubscriberId": {
          "reference": "DataSubscriber/86d46b60-6a86-4ead-83ef-ce7970d7c978"
        },
        "authenticationMethodIds": [
          {
            "reference": "AuthenticationMethod/8d89ee8f-86ed-4c59-83fa-fdd0fbf57161"
          }
        ]
      },
      "fullUrl": "ServiceAgent/8da12d84-4e4c-4c47-bdec-26b3614f6acf"
    }
  ],
  "link": [
    {
      "url": "ServiceAgent",
      "relation": "self"
    },
    {
      "url": "ServiceAgent?_page=1",
      "relation": "first"
    }
  ],
  "pageTotal": 2
}`)
		}
	})

	agents, resp, err := mdmClient.ServiceAgents.Get(nil)
	if !assert.Nilf(t, err, "unexpected error: %v", err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, agents) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, len(*agents))
}
