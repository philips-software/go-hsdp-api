package has_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/has"

	"github.com/stretchr/testify/assert"
)

func TestResourcesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	resourceID := "i-0f23dfed98e55913c"
	muxHAS.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resources": [
    {
      "id": "cke8gjrz10gew1305fzatnxm7",
      "resourceId": "i-08a214f0f844df43a",
      "organizationId": "`+hasOrgID+`",
      "imageId": "has-image-j4iikl0ie7b3",
      "resourceType": "g3s.xlarge",
      "clusterTag": "created-with-hs",
      "sessionId": "",
      "dns": "10.0.214.59",
      "state": "PENDING",
      "disabled": false,
      "region": "eu-west-1"
    }
  ]
}`)
		case "DELETE":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "results": [
    {
      "resourceId": "i-08a214f0f844df43a",
      "action": "DELETE",
      "resultCode": 200,
      "resultMessage": "Success"
    }
  ]
}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resources": [
    {
      "id": "cke8gjrz10gew1305fzatnxm7",
      "resourceId": "i-08a214f0f844df43a",
      "organizationId": "`+hasOrgID+`",
      "imageId": "has-image-j4iikl0ie7b3",
      "resourceType": "g3s.xlarge",
      "clusterTag": "created-with-hs",
      "sessionId": "",
      "dns": "10.0.214.59",
      "state": "RUNNING",
      "disabled": false,
      "region": "eu-west-1"
    }
  ]
}`)
		}
	})
	muxHAS.HandleFunc("/resource/stop", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "results": [
    {
      "resourceId": "i-0f23dfed98e55913c",
      "action": "STOP",
      "resultCode": 200,
      "resultMessage": "Success"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxHAS.HandleFunc("/resource/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "results": [
    {
      "resourceId": "i-0f23dfed98e55913c",
      "action": "START",
      "resultCode": 200,
      "resultMessage": "Success"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxHAS.HandleFunc("/resource/"+resourceID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
      "id": "cke8gjrz10gew1305fzatnxm7",
      "resourceId": "i-08a214f0f844df43a",
      "organizationId": "`+hasOrgID+`",
      "imageId": "has-image-j4iikl0ie7b3",
      "resourceType": "g3s.xlarge",
      "clusterTag": "created-with-hs",
      "sessionId": "",
      "dns": "10.0.214.59",
      "state": "RUNNING",
      "disabled": false,
      "region": "eu-west-1"
    }`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
      "id": "cke8gjrz10gew1305fzatnxm7",
      "resourceId": "i-08a214f0f844df43a",
      "organizationId": "`+hasOrgID+`",
      "imageId": "has-image-j4iikl0ie7b3",
      "resourceType": "g3s.xlarge",
      "clusterTag": "created-with-hs",
      "sessionId": "",
      "dns": "10.0.214.59",
      "state": "RUNNING",
      "disabled": false,
      "region": "eu-west-1"
    }`)
		case "DELETE":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "results": [
    {
      "resourceId": "i-08a214f0f844df43a",
      "action": "DELETE",
      "resultCode": 200,
      "resultMessage": "Success"
    }
  ]
}`)
		}
	})

	r := has.Resource{
		ImageID:      "has-image-xxx",
		ResourceType: "g3s.xlarge",
		Count:        1,
		ClusterTag:   "created-with-hs",
		EBS: has.EBS{
			DeleteOnTermination: true,
			Encrypted:           true,
			VolumeSize:          50,
			VolumeType:          "standard",
		},
	}

	resources, resp, err := hasClient.Resources.CreateResource(r)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resources) {
		return
	}
	if !assert.Equal(t, 1, len(*resources)) {
		return
	}
	assert.Equal(t, "has-image-j4iikl0ie7b3", (*resources)[0].ImageID)

	resource, resp, err := hasClient.Resources.GetResource(resourceID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resource) {
		return
	}
	assert.Equal(t, "has-image-j4iikl0ie7b3", resource.ImageID)

	report, resp, err := hasClient.Resources.StopResource([]string{resourceID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, report) {
		return
	}
	report, resp, err = hasClient.Resources.StartResource([]string{resourceID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, report) {
		return
	}

	resources, resp, err = hasClient.Resources.GetResources(&has.ResourceOptions{
		ResourceID: &resourceID,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resources) {
		return
	}

	resourceResponse, resp, err := hasClient.Resources.DeleteResource(resourceID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resourceResponse) {
		return
	}

	resourceResponse, resp, err = hasClient.Resources.DeleteResources(&has.ResourceOptions{
		ImageID: &resourceID,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resourceResponse) {
		return
	}
}
