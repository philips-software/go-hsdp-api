package has_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-software/go-hsdp-api/has"
)

func TestSessionsCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	sessionID := "cke8qn6hs0gjb1305088jt9w6"
	imageID := "has-image-j4jjkl0ie7b3"
	resourceID := "i-0f23dfed98e55913c"

	muxHAS.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "sessions": [
    {
      "sessionId": "`+sessionID+`",
      "sessionUrl": "https://some.url/session?token=xxx#console",
      "state": "AVAILABLE",
      "region": "eu-west-1",
      "resourceId": "`+resourceID+`",
      "userId": "`+userUUID+`",
      "sessionType": "USER"
    }
  ]
}`)
		}
	})

	muxHAS.HandleFunc("/user/"+userUUID+"/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "sessions": [
    {
      "sessionId": "`+sessionID+`",
      "sessionUrl": "",
      "state": "PENDING",
      "region": "eu-west-1",
      "resourceId": "`+resourceID+`",
      "userId": "`+userUUID+`",
      "sessionType": "USER"
    }
  ]
}`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "sessions": [
    {
      "sessionId": "`+sessionID+`",
      "sessionUrl": "https://some.url/session?token=xxx#console",
      "state": "AVAILABLE",
      "region": "eu-west-1",
      "resourceId": "`+resourceID+`",
      "userId": "`+userUUID+`",
      "sessionType": "USER"
    }
  ]
}`)
		}
	})

	sessions, resp, err := hasClient.Sessions.CreateSession(userUUID, has.Session{
		Region:     "eu-west-1",
		ImageID:    imageID,
		ClusterTag: "created-with-hs",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sessions) {
		return
	}

	sessions, resp, err = hasClient.Sessions.GetSessions()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sessions) {
		return
	}

	sessions, resp, err = hasClient.Sessions.GetSession(userUUID, &has.SessionOptions{
		ResourceID: &resourceID,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, sessions) {
		return
	}

	ok, resp, err := hasClient.Sessions.DeleteSession(userUUID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.True(t, ok) {
		return
	}
}
