package has_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImages(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxHAS.HandleFunc("/has/image", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "images": [
    {
      "id": "has-image-j4aakl0ie7b3",
      "name": "BaseImage2k12r2",
      "description": "Base image for Windows Server 2012 R2",
      "regions": [
        "us-east-1",
        "us-east-2",
        "us-west-2",
        "eu-west-1"
      ]
    },
    {
      "id": "has-image-98fur6fbsiff",
      "name": "BaseImageWS2k16",
      "description": "Base Image for windows server 2016",
      "regions": [
        "us-east-1"
      ]
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	images, resp, err := hasClient.Images.GetImages()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, images) {
		return
	}
	assert.Equal(t, 2, len(*images))
}
